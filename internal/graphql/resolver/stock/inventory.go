package stock

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// inventoryRef identifies one inventory annotation set of the stock
// center.
type inventoryRef struct {
	tag      string
	ontology string
}

// inventoryRefs maps each GraphQL stock kind to its inventory set.
var inventoryRefs = map[models.StockEntityType]inventoryRef{
	models.StockEntityTypeStrain: {
		tag:      registry.StrainInvTag,
		ontology: registry.StrainInvOnto,
	},
	models.StockEntityTypePlasmid: {
		tag:      registry.PlasmidInvTag,
		ontology: registry.PlasmidInvOnto,
	},
}

// inventoryRefFor resolves the inventory set of a stock kind. An
// unknown kind fails, because the stock center has no inventory set
// for it.
func inventoryRefFor(
	entity models.StockEntityType,
) E.Either[error, inventoryRef] {
	return F.Pipe2(
		inventoryRefs,
		R.Lookup[inventoryRef](entity),
		E.FromOption[inventoryRef](func() error {
			return fmt.Errorf(
				"unexpected stock entity %s for inventory lookup",
				entity,
			)
		}),
	)
}

var (
	// codeEq compares gRPC status codes.
	codeEq = EQ.FromStrictEquals[codes.Code]()

	// isNotFoundStatus is true when an error reports codes.NotFound.
	isNotFoundStatus = F.Flow2(
		status.Code,
		EQ.Equals(codeEq)(codes.NotFound),
	)

	// recoverNotFound turns a missing annotation into a not-in-stock
	// result and keeps every other error.
	recoverNotFound = P.Fold(
		func(err error) IOE.IOEither[error, bool] {
			return IOE.Left[bool](err)
		},
		func(error) IOE.IOEither[error, bool] {
			return IOE.Of[error](false)
		},
	)
)

// onInventoryError and onInventorySuccess fold one inventory lookup to
// the tuple shape of the stock resolver edges.
func onInventoryError(err error) T.Tuple2[error, bool] {
	return T.MakeTuple2(err, false)
}

func onInventorySuccess(inStock bool) T.Tuple2[error, bool] {
	return T.MakeTuple2[error](nil, inStock)
}

// inStockForEntity reports whether one stock is in the stock center
// inventory. The entity of the stock selects the inventory annotation
// set. A missing annotation gives false. Every other failure is an
// error.
func inStockForEntity(
	ctx context.Context,
	annoClient annotation.TaggedAnnotationServiceClient,
	entryID string,
	entity models.StockEntityType,
) IOE.IOEither[error, bool] {
	return F.Pipe2(
		inventoryRefFor(entity),
		IOE.FromEither[error, inventoryRef],
		IOE.Chain(func(ref inventoryRef) IOE.IOEither[error, bool] {
			return F.Pipe3(
				IOE.TryCatchError(func() (*annotation.TaggedAnnotation, error) {
					return annoClient.GetEntryAnnotation(
						ctx,
						&annotation.EntryAnnotationRequest{
							Tag:      ref.tag,
							Ontology: ref.ontology,
							EntryId:  entryID,
						},
					)
				}),
				IOE.MapTo[error, *annotation.TaggedAnnotation](true),
				IOE.OrElse(recoverNotFound(isNotFoundStatus)),
				IOE.MapLeft[bool](func(err error) error {
					return fmt.Errorf(
						"error getting %s from inventory: %w",
						entryID,
						err,
					)
				}),
			)
		}),
	)
}
