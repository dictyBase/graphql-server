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
	S "github.com/IBM/fp-go/v2/string"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

func PlasmidFilterToQuery(filter *models.PlasmidListFilter) (string, error) {
	if filter == nil {
		return "", nil
	}

	return E.UnwrapError(F.Pipe1(
		validatePlasmidFilter(filter),
		E.Map[error](func(_ struct{}) string {
			return buildPlasmidFieldQuery(filter)
		}),
	))
}

func validatePlasmidFilter(filter *models.PlasmidListFilter) E.Either[error, struct{}] {
	return E.MonadChain(
		validatePlasmidUnsupportedFields(filter),
		func(_ struct{}) E.Either[error, struct{}] {
			return validatePlasmidType(filter)
		},
	)
}

func validatePlasmidType(filter *models.PlasmidListFilter) E.Either[error, struct{}] {
	switch filter.PlasmidType {
	case models.PlasmidTypeAll:
		return E.Right[error](struct{}{})
	case models.PlasmidTypeRegular, models.PlasmidTypeGoldenBraid:
		return E.Left[struct{}](fmt.Errorf(
			"plasmid_type filter is not yet verified for stock query conversion",
		))
	}
	return E.Left[struct{}](fmt.Errorf("invalid plasmid type %s", filter.PlasmidType.String()))
}

func buildPlasmidFieldQuery(filter *models.PlasmidListFilter) string {
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

func validatePlasmidUnsupportedFields(filter *models.PlasmidListFilter) E.Either[error, struct{}] {
	if filter.ID != nil {
		return E.Left[struct{}](fmt.Errorf("id filter is not yet supported in stock query conversion"))
	}
	if filter.InStock != nil {
		return E.Left[struct{}](fmt.Errorf("in_stock filter is not yet supported in stock query conversion"))
	}
	return E.Right[error](struct{}{})
}
