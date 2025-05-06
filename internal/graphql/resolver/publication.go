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

// FetchPublicationDetailsParams defines the parameters for the fetchPublicationDetails function.
type FetchPublicationDetailsParams struct {
	Ctx        context.Context
	PubID      string
	GeneID     string // For logging context
	Qrs        *QueryResolver
	FeatClient feature.FeatureAnnotationServiceClient
}

// FetchGeneAnnotationParams defines the parameters for the fetchGeneAnnotation function.
type FetchGeneAnnotationParams struct {
	Ctx    context.Context
	Client feature.FeatureAnnotationServiceClient
	GeneID string
	Logger *logrus.Entry
}

// CollectPublicationResultsParams defines the parameters for the collectPublicationResults function.
type CollectPublicationResultsParams struct {
	Ctx     context.Context
	PubChan <-chan *models.PublicationWithGene
	ErrChan <-chan error
	Logger  *logrus.Entry
}

// LaunchPublicationFetchersParams defines the parameters for the launchPublicationFetchers method.
type LaunchPublicationFetchersParams struct {
	Ctx        context.Context
	Gene       string
	PubIDs     []string
	FeatClient feature.FeatureAnnotationServiceClient
	PubChan    chan<- *models.PublicationWithGene
	ErrChan    chan<- error
	Sem        *concurrency.Semaphore
	Wg         *sync.WaitGroup
	CancelFunc context.CancelFunc
	Qrs        *QueryResolver
}

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

	fetchDetailsParams := &FetchPublicationDetailsParams{
		Ctx:        params.ctx,
		PubID:      params.pubID,
		GeneID:     params.gene,
		Qrs:        params.qrs,
		FeatClient: params.featClient,
	}
	pubWithGene, err := fetchPublicationDetails(fetchDetailsParams)
	if err != nil {
		params.errChan <- fmt.Errorf("failed fetching details for pub %s: %w", params.pubID, err)
		params.cancelFunc() // Cancel other ongoing operations
		return
	}
	params.pubChan <- pubWithGene
}

// launchPublicationFetchers starts goroutines to fetch publication details concurrently.
func (qrs *QueryResolver) launchPublicationFetchers(
	params *LaunchPublicationFetchersParams,
) {
	for _, pubID := range params.PubIDs {
		params.Wg.Add(1)
		params.Sem.Acquire()
		fetchParams := &fetchPubAsyncParams{
			ctx:        params.Ctx,
			pubID:      pubID,
			gene:       params.Gene,
			featClient: params.FeatClient,
			pubChan:    params.PubChan,
			errChan:    params.ErrChan,
			sem:        params.Sem,
			wg:         params.Wg,
			cancelFunc: params.CancelFunc,
			qrs:        params.Qrs, // Pass the resolver
		}
		go qrs.fetchPublicationAsync(fetchParams)
	}
	go func() {
		params.Wg.Wait()
		close(params.PubChan)
		close(params.ErrChan)
	}()
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
	// 1. Fetch the initial gene annotation
	featClient := qrs.GetFeatAnnotationClient(
		registry.FeatAnno,
	)
	feat, err := fetchGeneAnnotation(&FetchGeneAnnotationParams{
		Ctx:    ctx,
		Client: featClient,
		GeneID: gene,
		Logger: qrs.Logger,
	})
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
	sem := concurrency.NewSemaphore(6) // Limit concurrency
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
	qrs.launchPublicationFetchers(&LaunchPublicationFetchersParams{
		Ctx:        fetchCtx,
		Gene:       gene,
		PubIDs:     pubIDs,
		FeatClient: featClient,
		PubChan:    pubChan,
		ErrChan:    errChan,
		Sem:        sem,
		Wg:         &sync.WaitGroup{},
		CancelFunc: cancelFetch,
		Qrs:        qrs,
	})
	// 6. Collect results and handle errors
	pubListResult, firstErr := collectPublicationResults(
		&CollectPublicationResultsParams{
			Ctx:     ctx,
			PubChan: pubChan,
			ErrChan: errChan,
			Logger:  qrs.Logger,
		},
	)
	// If any error occurred during fetching, return immediately
	if firstErr != nil {
		return nil, fmt.Errorf(
			"encountered errors while fetching publications: %w",
			firstErr,
		)
	}
	// 8. Return the collected list
	return pubListResult, nil
}

// collectPublicationResults waits for fetchers, collects results, and handles
// errors.
func collectPublicationResults(
	params *CollectPublicationResultsParams,
) ([]*models.PublicationWithGene, error) {
	pubList := make([]*models.PublicationWithGene, 0)
	var firstErr error

	// Process errors first to add them to GraphQL context
	for fetchErr := range params.ErrChan {
		errorutils.AddGQLError(params.Ctx, fetchErr) // Add specific fetch error
		params.Logger.Error(fetchErr)                // Log the specific error
		if firstErr == nil {
			firstErr = fetchErr // Capture the first error
		}
	}

	// Collect successful results
	for pub := range params.PubChan {
		pubList = append(pubList, pub)
	}

	return pubList, firstErr
}

// fetchGeneAnnotation retrieves the feature annotation for a given gene ID.
// It handles NotFound errors gracefully by returning nil, nil.
func fetchGeneAnnotation(
	params *FetchGeneAnnotationParams,
) (*feature.FeatureAnnotation, error) {
	feat, err := params.Client.GetFeatureAnnotation(
		params.Ctx,
		&feature.FeatureAnnotationId{Id: params.GeneID},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			params.Logger.Warnf("gene %s not found", params.GeneID)
			return nil, nil // Return nil, nil for NotFound
		}
		errorutils.AddGQLError(params.Ctx, err)
		params.Logger.Errorf(
			"error fetching feature annotation for gene ID %s: %v",
			params.GeneID,
			err,
		)
		return nil, fmt.Errorf(
			"error fetching feature annotation for gene %s: %w",
			params.GeneID,
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
	params *FetchPublicationDetailsParams,
) (*models.PublicationWithGene, error) {
	pub, err := fetch.FetchPublication(
		params.Ctx, params.Qrs.GetRedisRepository(cache.RedisKey),
		params.Qrs.Registry.GetAPIEndpoint(registry.PUBLICATION), params.PubID,
	)
	if err != nil {
		// Error already logged by FetchPublication if needed, just wrap and return
		return nil, fmt.Errorf(
			"error fetching publication %s: %w",
			params.PubID,
			err,
		)
	}
	pubWithGene := &models.PublicationWithGene{
		ID:       params.PubID, // Use params.PubID here, assuming pub.ID was a mistake in the diff
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
	featAnnos, err := params.FeatClient.ListFeatureAnnotationsByPubmedId(
		params.Ctx,
		&feature.PubmedId{Id: params.PubID}, // Corrected to use PublicationId
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			params.Qrs.Logger.Warnf(
				"genes for pubmed Id %s not found",
				params.PubID,
			)
			// return with empty relatedGenes
			return pubWithGene, nil
		}
		// Log warning but don't fail the whole request for this publication
		errMsg := fmt.Errorf(
			"error fetching feature annotations for pubmed ID %s (related to gene %s): %v",
			params.PubID,
			params.GeneID,
			err,
		)
		params.Qrs.Logger.Error(errMsg)
		return nil, errMsg
	}
	pubWithGene.RelatedGenes = collection.Map(
		featAnnos.Data,
		featureAnnotationToGene,
	)
	return pubWithGene, nil
}
