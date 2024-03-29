// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

func (ec *executionContext) unmarshalNTimestamp2timeᚐTime(ctx context.Context, v interface{}) (time.Time, error) {
	res, err := models.UnmarshalTimestamp(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNTimestamp2timeᚐTime(ctx context.Context, sel ast.SelectionSet, v time.Time) graphql.Marshaler {
	res := models.MarshalTimestamp(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalNTimestamp2ᚖtimeᚐTime(ctx context.Context, v interface{}) (*time.Time, error) {
	res, err := models.UnmarshalTimestamp(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNTimestamp2ᚖtimeᚐTime(ctx context.Context, sel ast.SelectionSet, v *time.Time) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	res := models.MarshalTimestamp(*v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalOTimestamp2ᚖtimeᚐTime(ctx context.Context, v interface{}) (*time.Time, error) {
	if v == nil {
		return nil, nil
	}
	res, err := models.UnmarshalTimestamp(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOTimestamp2ᚖtimeᚐTime(ctx context.Context, sel ast.SelectionSet, v *time.Time) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	res := models.MarshalTimestamp(*v)
	return res
}

// endregion ***************************** type.gotpl *****************************
