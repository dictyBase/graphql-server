package resolverutils

import (
	"fmt"
	"strings"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

func PlasmidFilterToQuery(filter *models.PlasmidListFilter) (string, error) {
	if filter == nil {
		return "", nil
	}

	var query strings.Builder

	fieldQuery, err := plasmidFieldsQuery(filter)
	if err != nil {
		return "", err
	}
	if fieldQuery != "" {
		query.WriteString(fieldQuery)
	}

	typeQuery, err := plasmidTypeQuery(filter)
	if err != nil {
		return "", err
	}
	if typeQuery != "" {
		if query.Len() > 0 {
			query.WriteString(";")
		}
		query.WriteString(typeQuery)
	}

	return query.String(), nil
}

func plasmidFieldsQuery(filter *models.PlasmidListFilter) (string, error) {
	var query strings.Builder

	if filter.Summary != nil {
		fmt.Fprintf(&query, "summary=~%s", *filter.Summary)
	}
	if filter.Name != nil {
		if query.Len() > 0 {
			query.WriteString(";")
		}
		fmt.Fprintf(&query, "plasmid_name===%s", *filter.Name)
	}
	if filter.ID != nil {
		return "", fmt.Errorf("id filter is not yet supported in stock query conversion")
	}
	if filter.InStock != nil {
		return "", fmt.Errorf("in_stock filter is not yet supported in stock query conversion")
	}

	return query.String(), nil
}

func plasmidTypeQuery(filter *models.PlasmidListFilter) (string, error) {
	switch filter.PlasmidType {
	case models.PlasmidTypeAll:
		return fmt.Sprintf(
			"ontology==%s;tag==%s,tag==%s",
			registry.DictyPlasmidPropOntology,
			registry.RegularPlasmidTag,
			registry.GoldenBraidPlasmidTag,
		), nil
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
	}

	return "", fmt.Errorf("invalid plasmid type %s", filter.PlasmidType.String())
}
