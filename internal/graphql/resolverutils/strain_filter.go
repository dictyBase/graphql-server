package resolverutils

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

// strainFilterQueryPair carries both the original filter and its validated
// strain type through the assembly pipeline.
type strainFilterQueryPair = T.Tuple2[*models.StrainListFilter, models.StrainType]

// ── validation predicates ────────────────────────────────────────────────────

var (
	isNilStrainID = F.Pipe1(
		P.IsZero[*string](),
		P.ContraMap(func(f *models.StrainListFilter) *string { return f.ID }),
	)
	CheckStrainIDField = E.FromPredicate(
		isNilStrainID,
		func(f *models.StrainListFilter) error {
			return fmt.Errorf(
				"strain list filter %v: id filter is not yet supported", f,
			)
		},
	)

	isNilStrainInStock = F.Pipe1(
		P.IsZero[*bool](),
		P.ContraMap(func(f *models.StrainListFilter) *bool { return f.InStock }),
	)
	CheckStrainInStockField = E.FromPredicate(
		isNilStrainInStock,
		func(f *models.StrainListFilter) error {
			return fmt.Errorf(
				"strain list filter %v: in_stock filter is not yet supported", f,
			)
		},
	)
)

// ── type → DSL lookup table ──────────────────────────────────────────────────

// strainTypeQueryMap maps each StrainType to its annotation-service filter DSL.
// BACTERIAL intentionally uses BacterialFoodSourceTag (the sole criterion).
// ALL intentionally excludes bacterial strains (separate category).
var strainTypeQueryMap = map[models.StrainType]string{
	models.StrainTypeAll: fmt.Sprintf(
		"ontology===%s;tag===%s,tag===%s",
		registry.DictyStrainPropOntology,
		registry.GeneralStrainTag,
		registry.GwdiStrainTag,
	),
	models.StrainTypeRegular: fmt.Sprintf(
		"ontology===%s;tag===%s",
		registry.DictyStrainPropOntology,
		registry.GeneralStrainTag,
	),
	models.StrainTypeGwdi: fmt.Sprintf(
		"ontology===%s;tag===%s",
		registry.DictyStrainPropOntology,
		registry.GwdiStrainTag,
	),
	models.StrainTypeBacterial: fmt.Sprintf(
		"ontology===%s;tag===%s",
		registry.StrainCharOnto,
		registry.BacterialFoodSourceTag,
	),
}

var strainTypeIsValid = E.FromPredicate(
	func(st models.StrainType) bool { return st.IsValid() },
	func(st models.StrainType) error {
		return fmt.Errorf("invalid strain type %s", st.String())
	},
)

// ── field query builder ──────────────────────────────────────────────────────

var formatStrainFieldQuery = F.Curry2(
	func(format string, value *string) string {
		return S.Format[string](format)(*value)
	},
)

var compactStrainOptionStrings = A.FilterMap(O.Fold(
	F.Constant(O.None[string]()),
	O.Some[string],
))

// ── exported API ─────────────────────────────────────────────────────────────

// BacterialAnnotationFilter returns the annotation service filter DSL that
// matches strains annotated as "bacterial food source" in
// dicty_strain_characteristics — the sole criterion for bacterial strains.
func BacterialAnnotationFilter() string {
	return fmt.Sprintf(
		"ontology===%s;tag===%s",
		registry.StrainCharOnto,
		registry.BacterialFoodSourceTag,
	)
}

// assembleStrainFilterQuery assembles the annotation-service filter DSL query
// string from the filter's label, summary, and strain type fields.
func assembleStrainFilterQuery(p strainFilterQueryPair) string {
	f := p.F1
	st := p.F2
	return F.Pipe2(
		[]O.Option[string]{
			F.Pipe1(
				O.FromNillable(f.Label),
				O.Map(formatStrainFieldQuery("label=~%s")),
			),
			F.Pipe1(
				O.FromNillable(f.Summary),
				O.Map(formatStrainFieldQuery("summary=~%s")),
			),
			F.Pipe1(
				R.Lookup[string](st)(strainTypeQueryMap),
				O.Chain(func(s string) O.Option[string] {
					return O.FromPredicate(
						func(v string) bool { return len(v) > 0 },
					)(s)
				}),
			),
		},
		compactStrainOptionStrings,
		A.Intercalate(S.Monoid)(";"),
	)
}

// StrainFilterToQueryFP converts a StrainListFilter to an annotation-service
// internal query string, validating the filter's fields and returning an error
// if any are invalid or unsupported. The filter's label and summary fields are
// converted to regex queries, and the strain type is used to look up the
// appropriate query string from the strainTypeQueryMap. The filter's id and
// in_stocking fields are not yet supported and will cause an error if
// present.
func StrainFilterToQueryFP(filter *models.StrainListFilter) IOE.IOEither[error, string] {
	return F.Pipe3(
		filter,
		O.FromNillable[models.StrainListFilter],
		O.Fold(
			F.Constant(E.Right[error]("")),
			func(f *models.StrainListFilter) E.Either[error, string] {
				return F.Pipe7(
					f,
					E.Of[error, *models.StrainListFilter],
					E.Chain(CheckStrainIDField),
					E.Chain(CheckStrainInStockField),
					E.Map[error](func(f *models.StrainListFilter) models.StrainType {
						return f.StrainType
					}),
					E.Chain(strainTypeIsValid),
					E.Map[error](func(st models.StrainType) strainFilterQueryPair {
						return T.MakeTuple2(f, st)
					}),
					E.Map[error](assembleStrainFilterQuery),
				)
			},
		),
		IOE.FromEither[error, string],
	)
}
