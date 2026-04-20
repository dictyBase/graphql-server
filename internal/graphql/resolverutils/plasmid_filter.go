// Package resolverutils provides utility functions for GraphQL resolvers,
// including functions to convert GraphQL filters into query strings for
// database queries.
package resolverutils

import (
	"fmt"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

var (
	isNilID = F.Pipe1(
		P.IsZero[*string](),
		P.ContraMap(func(f *models.PlasmidListFilter) *string {
			return f.ID
		}),
	)
	CheckIDField = E.FromPredicate(
		isNilID,
		func(filter *models.PlasmidListFilter) error {
			return fmt.Errorf(
				"plasmid list filter %v: id filter is not yet supported in stock query conversion",
				filter,
			)
		},
	)

	isNilInStock = F.Pipe1(
		P.IsZero[*bool](),
		P.ContraMap(func(f *models.PlasmidListFilter) *bool {
			return f.InStock
		}),
	)
	CheckInStockField = E.FromPredicate(
		isNilInStock,
		func(filter *models.PlasmidListFilter) error {
			return fmt.Errorf(
				"plasmid list filter %v: in_stock filter is not yet supported in stock query conversion",
				filter,
			)
		},
	)
)

func PlasmidFilterToQuery(filter *models.PlasmidListFilter) (string, error) {
	if filter == nil {
		return "", nil
	}

	if _, err := E.UnwrapError(F.Pipe2(
		filter,
		CheckIDField,
		E.Chain(CheckInStockField),
	)); err != nil {
		return "", err
	}

	fieldQuery := BuildPlasmidFieldQuery(filter)
	typeQuery, err := plasmidTypeQuery(filter)
	if err != nil {
		return fieldQuery, err
	}
	if fieldQuery == "" {
		return typeQuery, nil
	}
	if typeQuery == "" {
		return fieldQuery, nil
	}

	return strings.Join([]string{fieldQuery, typeQuery}, ";"), nil
}

func plasmidTypeQuery(filter *models.PlasmidListFilter) (string, error) {
	switch filter.PlasmidType {
	case models.PlasmidTypeAll:
		return "", nil
	case models.PlasmidTypeRegular:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyPlasmidPropOntology,
			registry.RegularPlasmidTag,
		), nil
	case models.PlasmidTypeGoldenBraid:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyPlasmidPropOntology,
			registry.GoldenBraidPlasmidTag,
		), nil
	default:
		return "", fmt.Errorf("invalid plasmid type %s", filter.PlasmidType.String())
	}
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
