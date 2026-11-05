package fetch

import (
	"testing"

	"github.com/dictyBase/fp-go-loom/loomcheck"
	"github.com/dictyBase/fp-go-loom/pipelinecheck"
)

func TestPipelineCheck(t *testing.T) {
	t.Parallel()

	cfg := pipelinecheck.Config{
		Roots:                     []string{"publication.go"},
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

func TestLoomCheck(t *testing.T) {
	t.Parallel()
	loomcheck.Require(t, loomcheck.Config{Roots: []string{"."}})
}
