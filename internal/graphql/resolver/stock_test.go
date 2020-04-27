package resolver

import (
	"context"
	"testing"

	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPlasmid(t *testing.T) {
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	p, err := q.Plasmid(context.Background(), "DBP123456")
	assert.NoError(err, "expect no error from getting plasmid by ID")
	assert.Exactly(p.Data.Id, "DBP123456", "should match plasmid ID")
	assert.Exactly(p.Data.Attributes.CreatedBy, "art@vandelay.com", "should match created_by")
	assert.Exactly(p.Data.Attributes.UpdatedBy, "art@vandelay.com", "should match updated_by")
	assert.Exactly(p.Data.Attributes.Summary, "test summary", "should match summary")
	assert.Exactly(p.Data.Attributes.EditableSummary, "editable test summary", "should match editable summary")
	assert.Exactly(p.Data.Attributes.Depositor, "Kenny Bania", "should match depositor (he's gold)")
	assert.ElementsMatch(p.Data.Attributes.Genes, []string{"sadA"}, "should match genes list")
	assert.ElementsMatch(p.Data.Attributes.Dbxrefs, []string{"test1"}, "should match dbxrefs")
	assert.ElementsMatch(p.Data.Attributes.Publications, []string{"99999"}, "should match publications")
	assert.Exactly(p.Data.Attributes.ImageMap, "https://eric.dictybase.dev/test.jpg", "should match image map")
	assert.Exactly(p.Data.Attributes.Sequence, "ABCDEF", "should match sequence")
	assert.Exactly(p.Data.Attributes.Name, "pTest", "should match name")
}

func TestStrain(t *testing.T) {
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	p, err := q.Strain(context.Background(), "DBS123456")
	assert.NoError(err, "expect no error from getting strain by ID")
	assert.Exactly(p.Data.Id, "DBS123456", "should match strain ID")
	assert.Exactly(p.Data.Attributes.CreatedBy, "art@vandelay.com", "should match created_by")
	assert.Exactly(p.Data.Attributes.UpdatedBy, "art@vandelay.com", "should match updated_by")
	assert.Exactly(p.Data.Attributes.Summary, "test summary", "should match summary")
	assert.Exactly(p.Data.Attributes.EditableSummary, "editable test summary", "should match editable summary")
	assert.Exactly(p.Data.Attributes.Depositor, "Kenny Bania", "should match depositor (he's gold)")
	assert.ElementsMatch(p.Data.Attributes.Genes, []string{"sadA"}, "should match genes list")
	assert.ElementsMatch(p.Data.Attributes.Dbxrefs, []string{"test1"}, "should match dbxrefs")
	assert.ElementsMatch(p.Data.Attributes.Publications, []string{"99999"}, "should match publications")
	assert.Exactly(p.Data.Attributes.Label, "test99", "should match label")
	assert.Exactly(p.Data.Attributes.Species, "human", "should match species")
	assert.Exactly(p.Data.Attributes.Plasmid, "pTest", "should match plasmid")
	assert.ElementsMatch(p.Data.Attributes.Names, []string{"fusilli"}, "should match names")
}
