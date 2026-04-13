package resolverutils

import (
	"fmt"
	"strings"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetCursorFP(cursor *int) int64 {
	return F.Pipe2(
		O.FromNillable(cursor),
		O.Map(func(ptr *int) int64 { return int64(*ptr) }),
		O.GetOrElse(F.Constant(int64(0))),
	)
}

func GetLimitFP(limit *int) int64 {
	return F.Pipe2(
		O.FromNillable(limit),
		O.Map(func(ptr *int) int64 { return int64(*ptr) }),
		O.GetOrElse(F.Constant(int64(10))),
	)
}

func GetCursor(c *int) int64 { return GetCursorFP(c) }

func GetLimit(l *int) int64 { return GetLimitFP(l) }

func GetFilter(f *string) string {
	if f == nil {
		return ""
	}
	return *f
}

func strainTypeQuery(filter *models.StrainListFilter) (string, error) {
	switch filter.StrainType {
	case models.StrainTypeAll:
		return fmt.Sprintf(
			"ontology==%s;tag==%s,tag==%s,tag==%s",
			registry.DictyStrainPropOntology,
			registry.GeneralStrainTag,
			registry.GwdiStrainTag,
			registry.BacterialStrainTag,
		), nil
	case models.StrainTypeBacterial:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyStrainPropOntology,
			registry.BacterialStrainTag,
		), nil
	case models.StrainTypeRegular:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyStrainPropOntology,
			registry.GeneralStrainTag,
		), nil
	case models.StrainTypeGwdi:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyStrainPropOntology,
			registry.GwdiStrainTag,
		), nil
	}

	return "", fmt.Errorf("invalid strain type %s", filter.StrainType.String())
}

func strainFieldsQuery(filter *models.StrainListFilter) string {
	var query strings.Builder
	if filter.Label != nil {
		if query.Len() > 0 {
			query.WriteString(";")
		}
		fmt.Fprintf(&query, "label=~%s", *filter.Label)
	}
	if filter.Summary != nil {
		if query.Len() > 0 {
			query.WriteString(";")
		}
		fmt.Fprintf(&query, "summary=~%s", *filter.Summary)
	}

	return query.String()
}

func StrainFilterToQuery(filter *models.StrainListFilter) (string, error) {
	var query strings.Builder
	query.WriteString(strainFieldsQuery(filter))
	typeQuery, err := strainTypeQuery(filter)
	if err != nil {
		return query.String(), err
	}
	if query.Len() > 0 {
		query.WriteString(";")
	}
	query.WriteString(typeQuery)

	return query.String(), nil
}

func GetOntology(onto string) string {
	var oname string
	switch onto {
	case "phenotype":
		oname = registry.PhenoOntology
	case "characteristic":
		oname = registry.StrainCharOnto
	case "strain_inventory":
		oname = registry.StrainInvOnto
	case "plasmid_inventory":
		oname = registry.PlasmidInvOnto
	default:
		oname = "invalid ontology"
	}
	return oname
}

func TimeWithPointer(pbt *timestamppb.Timestamp) *time.Time {
	tstmp := aphgrpc.ProtoTimeStamp(pbt)
	return &tstmp
}
