// Package resolverutils provides utility functions for GraphQL resolvers,
// including functions to convert GraphQL filters into query strings for
// database queries.
package resolverutils

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

var isNilID = F.Pipe1(
	P.IsZero[*string](),
	P.ContraMap(func(f *models.PlasmidListFilter) *string { return f.ID }),
)

var CheckIDField = E.FromPredicate(
	isNilID,
	func(filter *models.PlasmidListFilter) error {
		return fmt.Errorf(
			"plasmid list filter %v: id filter is not yet supported in stock query conversion",
			filter,
		)
	},
)

var isNilInStock = F.Pipe1(
	P.IsZero[*bool](),
	P.ContraMap(func(f *models.PlasmidListFilter) *bool { return f.InStock }),
)

var CheckInStockField = E.FromPredicate(
	isNilInStock,
	func(filter *models.PlasmidListFilter) error {
		return fmt.Errorf(
			"plasmid list filter %v: in_stock filter is not yet supported in stock query conversion",
			filter,
		)
	},
)

var isNotRegular = P.Not(
	P.IsStrictEqual[models.PlasmidType]()(models.PlasmidTypeRegular),
)

var isNotGoldenBraid = P.Not(
	P.IsStrictEqual[models.PlasmidType]()(models.PlasmidTypeGoldenBraid),
)

var isUnverifiedPlasmidType = F.Pipe2(
	isNotRegular,
	P.And(isNotGoldenBraid),
	P.ContraMap(func(f *models.PlasmidListFilter) models.PlasmidType {
		return f.PlasmidType
	}),
)

var CheckUnverifiedPlasmidType = E.FromPredicate(
	isUnverifiedPlasmidType,
	func(filter *models.PlasmidListFilter) error {
		return fmt.Errorf(
			"plasmid list filter %v: plasmid_type filter is not yet verified for stock query conversion",
			filter,
		)
	},
)

var isAllPlasmidType = F.Pipe1(
	P.IsStrictEqual[models.PlasmidType]()(models.PlasmidTypeAll),
	P.ContraMap(func(f *models.PlasmidListFilter) models.PlasmidType {
		return f.PlasmidType
	}),
)

var CheckValidPlasmidType = E.FromPredicate(
	isAllPlasmidType,
	func(f *models.PlasmidListFilter) error {
		return fmt.Errorf("invalid plasmid type %s", f.PlasmidType.String())
	},
)

var validateFilterPipeline = F.Flow5(
	CheckIDField,
	E.Chain(CheckInStockField),
	E.Chain(CheckUnverifiedPlasmidType),
	E.Chain(CheckValidPlasmidType),
	E.Map[error](BuildPlasmidFieldQuery),
)

func PlasmidFilterToQuery(filter *models.PlasmidListFilter) (string, error) {
	return E.UnwrapError(F.Pipe1(
		O.FromNillable(filter),
		O.Fold(
			F.Constant(E.Right[error]("")),
			validateFilterPipeline,
		),
	))
}

func BuildPlasmidFieldQuery(filter *models.PlasmidListFilter) string {
	return F.Pipe2(
		[]O.Option[string]{
			F.Pipe1(
				O.FromNillable(filter.Summary),
				O.Map(func(ptr *string) string {
					return S.Format[string]("summary=~%s")(*ptr)
				}),
			),
			F.Pipe1(
				O.FromNillable(filter.Name),
				O.Map(func(ptr *string) string {
					return S.Format[string]("plasmid_name===%s")(*ptr)
				}),
			),
		},
		A.FilterMap(O.Fold(
			F.Constant(O.None[string]()),
			O.Some[string],
		)),
		A.Intercalate(S.Monoid)(";"),
	)
}
