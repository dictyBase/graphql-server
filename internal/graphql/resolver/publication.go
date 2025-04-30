package resolver

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/graphql-server/internal/collection"
	"github.com/dictyBase/graphql-server/internal/concurrency"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
)

// fetchPubAsyncParams holds the parameters for the fetchPublicationAsync goroutine.
type fetchPubAsyncParams struct {
	ctx        context.Context
	pubID      string
	gene       string
	featClient feature.FeatureAnnotationServiceClient
	pubChan    chan<- *models.PublicationWithGene
	errChan    chan<- error
	sem        *concurrency.Semaphore
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
	qrs        *QueryResolver
}

// Publication is the resolver for getting an individual publication by ID.
func (qrs *QueryResolver) Publication(
	ctx context.Context,
	id string,
) (*models.Publication, error) {
	pub, err := fetch.FetchPublication(
		ctx, qrs.GetRedisRepository(cache.RedisKey),
		qrs.Registry.GetAPIEndpoint(registry.PUBLICATION), id,
	)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, fmt.Errorf(
			"error in fetching publication %s",
			err,
		)
	}
	return pub, nil
}

// AllPublications is the resolver for the allPublications field.
func (qrs *QueryResolver) AllPublications(
	ctx context.Context,
	gene string,
	limit *int,
	sortBy *string,
) (*models.NumberOfPublicationsWithGene, error) {
	return &models.NumberOfPublicationsWithGene{NumPubs: 0}, nil
}

// ListRecentPublications is the resolver for the listRecentPublications field.
func (qrs *QueryResolver) ListRecentPublications(
	ctx context.Context,
	limit int,
) ([]*models.Publication, error) {
	return []*models.Publication{}, nil
}

// fetchPublicationAsync is a helper method to fetch publication details asynchronously
func (qrs *QueryResolver) fetchPublicationAsync(params *fetchPubAsyncParams) {
	defer params.sem.Release()
	defer params.wg.Done()

	pubWithGene, err := fetchPublicationDetails(
		params.ctx,
		params.pubID,
		params.gene,
		params.qrs, // Pass the resolver itself
		params.featClient,
	)
	if err != nil {
		params.errChan <- fmt.Errorf("failed fetching details for pub %s: %w", params.pubID, err)
		params.cancelFunc() // Cancel other ongoing operations
		return
	}
	params.pubChan <- pubWithGene
}

// ListPublicationsWithGene fetches all publications associated with a gene ID
// and includes related genes for each publication. Returns an empty list if
// gene is not found or has no publications. If any error occurs during fetching
// of publication details (other than not-found errors), the entire operation
// fails and no results are returned.
func (qrs *QueryResolver) ListPublicationsWithGene(
	ctx context.Context,
	gene string,
) ([]*models.PublicationWithGene, error) {
	pubList := make([]*models.PublicationWithGene, 0)
	featClient := qrs.GetFeatAnnotationClient(registry.FEAT_ANNO)

	// 1. Fetch the initial gene annotation
	feat, err := fetchGeneAnnotation(ctx, featClient, gene, qrs.Logger)
	if err != nil {
		return nil, err // Error already logged and added to GraphQL errors
	}
	if feat == nil { // Gene not found
		return pubList, nil
	}

	// 2. Check for PubMed IDs
	pubIDs := feat.Attributes.Pubmed
	if len(pubIDs) == 0 {
		qrs.Logger.Warnf("no pubmed IDs found for gene %s", gene)
		return pubList, nil
	}

	// 3. Setup concurrency
	sem := concurrency.NewSemaphore(3) // Limit concurrency
	var wg sync.WaitGroup
	// Buffered channel to prevent goroutines from blocking indefinitely
	pubChan := make(chan *models.PublicationWithGene, len(pubIDs))
	errChan := make(
		chan error,
		len(pubIDs),
	) // Channel for errors from goroutines

	// Create a cancellable context to abort other operations on first error
	fetchCtx, cancelFetch := context.WithCancel(ctx)
	defer cancelFetch()

	// 4. Launch goroutines to fetch publications concurrently
	for _, pubID := range pubIDs {
		wg.Add(1)
		sem.Acquire()
		params := &fetchPubAsyncParams{
			ctx:        fetchCtx,
			pubID:      pubID,
			gene:       gene,
			featClient: featClient,
			pubChan:    pubChan,
			errChan:    errChan,
			sem:        sem,
			wg:         &wg,
			cancelFunc: cancelFetch,
			qrs:        qrs, // Pass the resolver
		}
		go qrs.fetchPublicationAsync(params)
	}

	// 5. Wait for all goroutines to complete and close channels
	go func() {
		wg.Wait()
		close(pubChan)
		close(errChan)
	}()

	// 6. Collect results and handle errors
	// Process errors first to add them to GraphQL context
	var firstErr error // Variable to store the first encountered error
	for fetchErr := range errChan {
		errorutils.AddGQLError(ctx, fetchErr) // Add specific fetch error
		qrs.Logger.Error(fetchErr)            // Log the specific error
		if firstErr == nil {
			firstErr = fetchErr // Capture the first error
		}
	}

	// If any error occurred during fetching, return immediately
	if firstErr != nil {
		return nil, fmt.Errorf(
			"encountered errors while fetching publications: %w",
			firstErr,
		)
	}

	// Collect successful results
	for pub := range pubChan {
		pubList = append(pubList, pub)
	}

	// 8. Return the collected list (potentially partial if fetches failed)
	return pubList, nil
}

// fetchGeneAnnotation retrieves the feature annotation for a given gene ID.
// It handles NotFound errors gracefully by returning nil, nil.
func fetchGeneAnnotation(
	ctx context.Context,
	client feature.FeatureAnnotationServiceClient,
	geneID string,
	logger *logrus.Entry,
) (*feature.FeatureAnnotation, error) {
	feat, err := client.GetFeatureAnnotation(
		ctx,
		&feature.FeatureAnnotationId{Id: geneID},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logger.Warnf("gene %s not found", geneID)
			return nil, nil // Return nil, nil for NotFound
		}
		errorutils.AddGQLError(ctx, err)
		logger.Errorf(
			"error fetching feature annotation for gene ID %s: %v",
			geneID,
			err,
		)
		return nil, fmt.Errorf(
			"error fetching feature annotation for gene %s: %w",
			geneID,
			err,
		)
	}
	return feat, nil
}

// featureAnnotationToGene converts a FeatureAnnotation protobuf message
// to a GraphQL Gene model.
func featureAnnotationToGene(
	faData *feature.FeatureAnnotation,
) *models.Gene {
	return &models.Gene{
		ID:   faData.Id,
		Name: faData.Attributes.Name,
	}
}

// fetchPublicationDetails fetches a single publication and its related genes.
func fetchPublicationDetails(
	ctx context.Context,
	pubID string,
	geneID string, // For logging context
	qrs *QueryResolver,
	featClient feature.FeatureAnnotationServiceClient,
) (*models.PublicationWithGene, error) {
	pub, err := fetch.FetchPublication(
		ctx, qrs.GetRedisRepository(cache.RedisKey),
		qrs.Registry.GetAPIEndpoint(registry.PUBLICATION), pubID,
	)
	if err != nil {
		// Error already logged by FetchPublication if needed, just wrap and return
		return nil, fmt.Errorf("error fetching publication %s: %w", pubID, err)
	}
	pubWithGene := &models.PublicationWithGene{
		ID:       pub.ID,
		Doi:      pub.Doi,
		Title:    pub.Title,
		Abstract: pub.Abstract,
		Journal:  pub.Journal,
		PubDate:  pub.PubDate,
		Volume:   pub.Volume,
		Pages:    pub.Pages,
		Issn:     pub.Issn,
		PubType:  pub.PubType,
		Source:   pub.Source,
		Issue:    pub.Issue,
		Status:   pub.Status,
		Authors:  pub.Authors,
	}

	// Fetch related genes for this publication
	featAnnos, err := featClient.ListFeatureAnnotationsByPubmedId(
		ctx,
		&feature.PubmedId{Id: pubID}, // Corrected to use PublicationId
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			qrs.Logger.Warnf("genes for pubmed Id %s not found", pubID)
			// return with empty relatedGenes
			return pubWithGene, nil
		}
		// Log warning but don't fail the whole request for this publication
		errMsg := fmt.Errorf(
			"error fetching feature annotations for pubmed ID %s (related to gene %s): %v",
			pubID,
			geneID,
			err,
		)
		qrs.Logger.Error(errMsg)
		return nil, errMsg
	}
	pubWithGene.RelatedGenes = collection.Map(
		featAnnos.Data,
		featureAnnotationToGene,
	)
	return pubWithGene, nil
}
