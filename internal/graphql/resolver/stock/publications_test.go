package stock

import (
	"bytes"
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/dictyBase/fp-go-loom/loomcheck"
	"github.com/dictyBase/fp-go-loom/pipelinecheck"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func newBufferLogger() (*logrus.Entry, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	logger := logrus.New()
	logger.SetOutput(buf)
	logger.SetLevel(logrus.DebugLevel)
	return logrus.NewEntry(logger), buf
}

func TestPublicationsPipelineCheck(t *testing.T) {
	t.Parallel()
	cfg := pipelinecheck.Config{
		Roots: []string{"publications.go"},
		IsEntrypoint: func(name string) bool {
			return name == "fetchStockPublications" || name == "getPhenotypes"
		},
		RequirePointFreeSeed:      true,
		RequirePointFreeBranching: true,
		RequireTrivialDelegate:    true,
	}
	pipelinecheck.Require(t, cfg)
	pipelinecheck.RequireNoNonRawTryCatchCallback(t, cfg)
	pipelinecheck.RequireNoIfErrInTryCatch(t, cfg)
	pipelinecheck.RequireBareIf(t, cfg)
	pipelinecheck.RequireTrivialDelegate(t, cfg)
	pipelinecheck.RequireNoImperativeLoop(t, cfg)
	pipelinecheck.RequireNoImperativeBranch(t, cfg)
	pipelinecheck.RequireNoDirectApply(t, cfg)
	pipelinecheck.RequireNoPopulatedStruct(t, cfg)
	pipelinecheck.RequireNoLocalPipeClosure(t, cfg)
	pipelinecheck.RequireNoRedundantFold(t, cfg)
	pipelinecheck.RequireSafePrintf(t, cfg)
	pipelinecheck.RequireNoHandRolledFileJoinInArrayMap(t, cfg)
}

func TestPublicationsLoomCheck(t *testing.T) {
	t.Parallel()
	loomcheck.Require(t, loomcheck.Config{
		Roots:               []string{"publications.go"},
		StrictAllPredicates: true,
	})
}

func TestFetchStockPublicationsSkipsInvalidIDs(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	logger, buf := newBufferLogger()

	res := F.Pipe5(
		stockPubsState{},
		stockPubsCtxLens.Set(context.Background()),
		stockPubsIDLens.Set("DBS0305923"),
		stockPubsLoggerLens.Set(logger),
		stockPubsRawIDsLens.Set([]string{"", "   ", "abc", "0"}),
		fetchStockPublications,
	)
	req.NoError(res.F1)
	req.NotNil(res.F2)
	req.Empty(res.F2)
	req.Contains(buf.String(), "skipping unresolvable publication identifiers")
	req.Contains(buf.String(), "DBS0305923")
}

func TestStrainAndPlasmidPublicationsSkipEmptyPMID(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	sr := strainResolver(mocks.MockedAnnotationClient())
	pr := plasmidResolver(mocks.MockedAnnotationClient())

	strainPubs, err := sr.Publications(
		context.Background(),
		&models.Strain{
			ID:        "DBS0305923",
			LazyStock: models.LazyStock{Publications: []string{""}},
		},
	)
	req.NoError(err)
	req.NotNil(strainPubs)
	req.Empty(strainPubs)

	plasmidPubs, err := pr.Publications(
		context.Background(),
		&models.Plasmid{
			ID:        "DBP0001102",
			LazyStock: models.LazyStock{Publications: []string{""}},
		},
	)
	req.NoError(err)
	req.NotNil(plasmidPubs)
	req.Empty(plasmidPubs)
}

func TestFilterResolvableIDs(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	req.Equal(
		[]string{"12062087", "10.1016/s1534-5807(02)00186-7"},
		filterResolvableIDs([]string{
			"",
			"abc",
			"12062087",
			"0",
			"10.1016/s1534-5807(02)00186-7",
		}),
	)
}

func TestGetPhenotypesSkipsUnresolvableLiteratureID(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	logger, buf := newBufferLogger()
	sr := strainResolver(mocks.MockedPhenotypeClient())
	sr.Logger = logger

	group := &annotation.TaggedAnnotationGroupCollection_Data{
		Group: &annotation.TaggedAnnotationGroup{
			Data: []*annotation.TaggedAnnotationGroup_Data{
				{
					Attributes: &annotation.TaggedAnnotationAttributes{
						Ontology: registry.PhenoOntology,
						Tag:      "decreased growth rate",
					},
				},
				{
					Attributes: &annotation.TaggedAnnotationAttributes{
						Ontology: registry.DictyAnnoOntology,
						Tag:      registry.LiteratureTag,
						Value:    "",
					},
				},
			},
		},
	}

	phenos := getPhenotypes(
		context.Background(),
		sr,
		[]*annotation.TaggedAnnotationGroupCollection_Data{group},
	)
	req.Len(phenos, 1)
	req.Equal("decreased growth rate", phenos[0].Phenotype)
	req.Nil(phenos[0].Publication)
	req.Contains(
		buf.String(),
		"skipping unresolvable phenotype publication identifier",
	)
}
