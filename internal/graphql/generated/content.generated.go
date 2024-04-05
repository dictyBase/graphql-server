// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

type ContentResolver interface {
	ID(ctx context.Context, obj *content.Content) (string, error)
	Name(ctx context.Context, obj *content.Content) (string, error)
	Slug(ctx context.Context, obj *content.Content) (string, error)
	CreatedBy(ctx context.Context, obj *content.Content) (*user.User, error)
	UpdatedBy(ctx context.Context, obj *content.Content) (*user.User, error)
	CreatedAt(ctx context.Context, obj *content.Content) (*time.Time, error)
	UpdatedAt(ctx context.Context, obj *content.Content) (*time.Time, error)
	Content(ctx context.Context, obj *content.Content) (string, error)
	Namespace(ctx context.Context, obj *content.Content) (string, error)
}

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Content_id(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_id(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().ID(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNID2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_id(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_name(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().Name(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_name(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_slug(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_slug(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().Slug(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_slug(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_created_by(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_created_by(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().CreatedBy(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*user.User)
	fc.Result = res
	return ec.marshalNUser2ᚖgithubᚗcomᚋdictyBaseᚋgoᚑgenprotoᚋdictybaseapisᚋuserᚐUser(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_created_by(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_User_id(ctx, field)
			case "first_name":
				return ec.fieldContext_User_first_name(ctx, field)
			case "last_name":
				return ec.fieldContext_User_last_name(ctx, field)
			case "email":
				return ec.fieldContext_User_email(ctx, field)
			case "organization":
				return ec.fieldContext_User_organization(ctx, field)
			case "group_name":
				return ec.fieldContext_User_group_name(ctx, field)
			case "first_address":
				return ec.fieldContext_User_first_address(ctx, field)
			case "second_address":
				return ec.fieldContext_User_second_address(ctx, field)
			case "city":
				return ec.fieldContext_User_city(ctx, field)
			case "state":
				return ec.fieldContext_User_state(ctx, field)
			case "zipcode":
				return ec.fieldContext_User_zipcode(ctx, field)
			case "country":
				return ec.fieldContext_User_country(ctx, field)
			case "phone":
				return ec.fieldContext_User_phone(ctx, field)
			case "is_active":
				return ec.fieldContext_User_is_active(ctx, field)
			case "created_at":
				return ec.fieldContext_User_created_at(ctx, field)
			case "updated_at":
				return ec.fieldContext_User_updated_at(ctx, field)
			case "roles":
				return ec.fieldContext_User_roles(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type User", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_updated_by(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_updated_by(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().UpdatedBy(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*user.User)
	fc.Result = res
	return ec.marshalNUser2ᚖgithubᚗcomᚋdictyBaseᚋgoᚑgenprotoᚋdictybaseapisᚋuserᚐUser(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_updated_by(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_User_id(ctx, field)
			case "first_name":
				return ec.fieldContext_User_first_name(ctx, field)
			case "last_name":
				return ec.fieldContext_User_last_name(ctx, field)
			case "email":
				return ec.fieldContext_User_email(ctx, field)
			case "organization":
				return ec.fieldContext_User_organization(ctx, field)
			case "group_name":
				return ec.fieldContext_User_group_name(ctx, field)
			case "first_address":
				return ec.fieldContext_User_first_address(ctx, field)
			case "second_address":
				return ec.fieldContext_User_second_address(ctx, field)
			case "city":
				return ec.fieldContext_User_city(ctx, field)
			case "state":
				return ec.fieldContext_User_state(ctx, field)
			case "zipcode":
				return ec.fieldContext_User_zipcode(ctx, field)
			case "country":
				return ec.fieldContext_User_country(ctx, field)
			case "phone":
				return ec.fieldContext_User_phone(ctx, field)
			case "is_active":
				return ec.fieldContext_User_is_active(ctx, field)
			case "created_at":
				return ec.fieldContext_User_created_at(ctx, field)
			case "updated_at":
				return ec.fieldContext_User_updated_at(ctx, field)
			case "roles":
				return ec.fieldContext_User_roles(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type User", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_created_at(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_created_at(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().CreatedAt(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*time.Time)
	fc.Result = res
	return ec.marshalNTimestamp2ᚖtimeᚐTime(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_created_at(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Timestamp does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_updated_at(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_updated_at(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().UpdatedAt(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*time.Time)
	fc.Result = res
	return ec.marshalNTimestamp2ᚖtimeᚐTime(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_updated_at(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Timestamp does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_content(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_content(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().Content(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_content(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Content_namespace(ctx context.Context, field graphql.CollectedField, obj *content.Content) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Content_namespace(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Content().Namespace(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Content_namespace(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Content",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _DeleteContent_success(ctx context.Context, field graphql.CollectedField, obj *models.DeleteContent) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_DeleteContent_success(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Success, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_DeleteContent_success(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "DeleteContent",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

func (ec *executionContext) unmarshalInputCreateContentInput(ctx context.Context, obj interface{}) (models.CreateContentInput, error) {
	var it models.CreateContentInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"name", "created_by", "content", "namespace"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "name":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("name"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Name = data
		case "created_by":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("created_by"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.CreatedBy = data
		case "content":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("content"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Content = data
		case "namespace":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("namespace"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Namespace = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUpdateContentInput(ctx context.Context, obj interface{}) (models.UpdateContentInput, error) {
	var it models.UpdateContentInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"id", "updated_by", "content"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "id":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
			data, err := ec.unmarshalNID2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.ID = data
		case "updated_by":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("updated_by"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.UpdatedBy = data
		case "content":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("content"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Content = data
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var contentImplementors = []string{"Content"}

func (ec *executionContext) _Content(ctx context.Context, sel ast.SelectionSet, obj *content.Content) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, contentImplementors)
	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Content")
		case "id":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_id(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "name":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_name(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "slug":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_slug(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "created_by":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_created_by(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "updated_by":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_updated_by(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "created_at":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_created_at(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "updated_at":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_updated_at(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "content":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_content(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		case "namespace":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Content_namespace(ctx, field, obj)
				if res == graphql.Null {
					atomic.AddUint32(&invalids, 1)
				}
				return res
			}

			out.Concurrently(i, func() graphql.Marshaler {
				return innerFunc(ctx)

			})
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch()
	if invalids > 0 {
		return graphql.Null
	}
	return out
}

var deleteContentImplementors = []string{"DeleteContent"}

func (ec *executionContext) _DeleteContent(ctx context.Context, sel ast.SelectionSet, obj *models.DeleteContent) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, deleteContentImplementors)
	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("DeleteContent")
		case "success":

			out.Values[i] = ec._DeleteContent_success(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch()
	if invalids > 0 {
		return graphql.Null
	}
	return out
}

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

func (ec *executionContext) marshalOContent2ᚖgithubᚗcomᚋdictyBaseᚋgoᚑgenprotoᚋdictybaseapisᚋcontentᚐContent(ctx context.Context, sel ast.SelectionSet, v *content.Content) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._Content(ctx, sel, v)
}

func (ec *executionContext) unmarshalOCreateContentInput2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐCreateContentInput(ctx context.Context, v interface{}) (*models.CreateContentInput, error) {
	if v == nil {
		return nil, nil
	}
	res, err := ec.unmarshalInputCreateContentInput(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalODeleteContent2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐDeleteContent(ctx context.Context, sel ast.SelectionSet, v *models.DeleteContent) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._DeleteContent(ctx, sel, v)
}

func (ec *executionContext) unmarshalOUpdateContentInput2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐUpdateContentInput(ctx context.Context, v interface{}) (*models.UpdateContentInput, error) {
	if v == nil {
		return nil, nil
	}
	res, err := ec.unmarshalInputUpdateContentInput(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

// endregion ***************************** type.gotpl *****************************
