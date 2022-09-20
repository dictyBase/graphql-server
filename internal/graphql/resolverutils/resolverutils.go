package resolverutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetCursor(c *int) int64 {
	if c == nil {
		return int64(0)
	}
	return int64(*c)
}

func GetLimit(l *int) int64 {
	if l == nil {
		return int64(10)
	}
	return int64(*l)
}

func GetFilter(f *string) string {
	if f == nil {
		return ""
	}
	return *f
}

func strainTypeQuery(filter *models.StrainListFilter) (string, error) {
	switch filter.StrainType {
	case models.StrainTypeEnumAll:
		return fmt.Sprintf(
			"ontology==%s;tag==%s,tag==%s,tag==%s",
			registry.DictyStrainPropOntology,
			registry.GeneralStrainTag,
			registry.GwdiStrainTag,
			registry.BacterialStrainTag,
		), nil
	case models.StrainTypeEnumBacterial:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyStrainPropOntology,
			registry.BacterialStrainTag,
		), nil
	case models.StrainTypeEnumRegular:
		return fmt.Sprintf(
			"ontology==%s;tag==%s",
			registry.DictyStrainPropOntology,
			registry.GeneralStrainTag,
		), nil
	case models.StrainTypeEnumGwdi:
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
		query.WriteString(fmt.Sprintf("label==%s", *filter.Label))
	}
	if filter.Summary != nil {
		if query.Len() > 0 {
			query.WriteString(fmt.Sprintf(";summary==%s", *filter.Summary))
		} else {
			query.WriteString(fmt.Sprintf("summary==%s", *filter.Summary))
		}
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
		query.WriteString(fmt.Sprintf(";%s", typeQuery))
	} else {
		query.WriteString(typeQuery)
	}

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
