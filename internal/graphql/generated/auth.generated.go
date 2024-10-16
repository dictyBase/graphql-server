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
	"github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

type AuthResolver interface {
	Identity(ctx context.Context, obj *auth.Auth) (*models.Identity, error)
}

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Auth_token(ctx context.Context, field graphql.CollectedField, obj *auth.Auth) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Auth_token(ctx, field)
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
		return obj.Token, nil
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

func (ec *executionContext) fieldContext_Auth_token(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Auth",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Auth_user(ctx context.Context, field graphql.CollectedField, obj *auth.Auth) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Auth_user(ctx, field)
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
		return obj.User, nil
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

func (ec *executionContext) fieldContext_Auth_user(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Auth",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
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

func (ec *executionContext) _Auth_identity(ctx context.Context, field graphql.CollectedField, obj *auth.Auth) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Auth_identity(ctx, field)
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
		return ec.resolvers.Auth().Identity(rctx, obj)
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
	res := resTmp.(*models.Identity)
	fc.Result = res
	return ec.marshalNIdentity2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐIdentity(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Auth_identity(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Auth",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Identity_id(ctx, field)
			case "identifier":
				return ec.fieldContext_Identity_identifier(ctx, field)
			case "provider":
				return ec.fieldContext_Identity_provider(ctx, field)
			case "user_id":
				return ec.fieldContext_Identity_user_id(ctx, field)
			case "created_at":
				return ec.fieldContext_Identity_created_at(ctx, field)
			case "updated_at":
				return ec.fieldContext_Identity_updated_at(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Identity", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_id(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_id(ctx, field)
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
		return obj.ID, nil
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

func (ec *executionContext) fieldContext_Identity_id(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_identifier(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_identifier(ctx, field)
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
		return obj.Identifier, nil
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

func (ec *executionContext) fieldContext_Identity_identifier(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_provider(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_provider(ctx, field)
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
		return obj.Provider, nil
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

func (ec *executionContext) fieldContext_Identity_provider(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_user_id(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_user_id(ctx, field)
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
		return obj.UserID, nil
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

func (ec *executionContext) fieldContext_Identity_user_id(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_created_at(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_created_at(ctx, field)
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
		return obj.CreatedAt, nil
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
	res := resTmp.(time.Time)
	fc.Result = res
	return ec.marshalNTimestamp2timeᚐTime(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Identity_created_at(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Timestamp does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Identity_updated_at(ctx context.Context, field graphql.CollectedField, obj *models.Identity) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Identity_updated_at(ctx, field)
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
		return obj.UpdatedAt, nil
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
	res := resTmp.(time.Time)
	fc.Result = res
	return ec.marshalNTimestamp2timeᚐTime(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Identity_updated_at(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Identity",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Timestamp does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Logout_success(ctx context.Context, field graphql.CollectedField, obj *models.Logout) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Logout_success(ctx, field)
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

func (ec *executionContext) fieldContext_Logout_success(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Logout",
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

func (ec *executionContext) unmarshalInputLoginInput(ctx context.Context, obj interface{}) (models.LoginInput, error) {
	var it models.LoginInput
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"client_id", "state", "code", "scopes", "provider", "redirect_url"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "client_id":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("client_id"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.ClientID = data
		case "state":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("state"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.State = data
		case "code":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("code"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Code = data
		case "scopes":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("scopes"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Scopes = data
		case "provider":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("provider"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Provider = data
		case "redirect_url":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("redirect_url"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.RedirectURL = data
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var authImplementors = []string{"Auth"}

func (ec *executionContext) _Auth(ctx context.Context, sel ast.SelectionSet, obj *auth.Auth) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, authImplementors)
	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Auth")
		case "token":

			out.Values[i] = ec._Auth_token(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&invalids, 1)
			}
		case "user":

			out.Values[i] = ec._Auth_user(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&invalids, 1)
			}
		case "identity":
			field := field

			innerFunc := func(ctx context.Context) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Auth_identity(ctx, field, obj)
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

var identityImplementors = []string{"Identity"}

func (ec *executionContext) _Identity(ctx context.Context, sel ast.SelectionSet, obj *models.Identity) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, identityImplementors)
	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Identity")
		case "id":

			out.Values[i] = ec._Identity_id(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "identifier":

			out.Values[i] = ec._Identity_identifier(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "provider":

			out.Values[i] = ec._Identity_provider(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "user_id":

			out.Values[i] = ec._Identity_user_id(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "created_at":

			out.Values[i] = ec._Identity_created_at(ctx, field, obj)

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "updated_at":

			out.Values[i] = ec._Identity_updated_at(ctx, field, obj)

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

var logoutImplementors = []string{"Logout"}

func (ec *executionContext) _Logout(ctx context.Context, sel ast.SelectionSet, obj *models.Logout) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, logoutImplementors)
	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Logout")
		case "success":

			out.Values[i] = ec._Logout_success(ctx, field, obj)

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

func (ec *executionContext) marshalNIdentity2githubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐIdentity(ctx context.Context, sel ast.SelectionSet, v models.Identity) graphql.Marshaler {
	return ec._Identity(ctx, sel, &v)
}

func (ec *executionContext) marshalNIdentity2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐIdentity(ctx context.Context, sel ast.SelectionSet, v *models.Identity) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec._Identity(ctx, sel, v)
}

func (ec *executionContext) marshalOAuth2ᚖgithubᚗcomᚋdictyBaseᚋgoᚑgenprotoᚋdictybaseapisᚋauthᚐAuth(ctx context.Context, sel ast.SelectionSet, v *auth.Auth) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._Auth(ctx, sel, v)
}

func (ec *executionContext) unmarshalOLoginInput2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐLoginInput(ctx context.Context, v interface{}) (*models.LoginInput, error) {
	if v == nil {
		return nil, nil
	}
	res, err := ec.unmarshalInputLoginInput(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOLogout2ᚖgithubᚗcomᚋdictyBaseᚋgraphqlᚑserverᚋinternalᚋgraphqlᚋmodelsᚐLogout(ctx context.Context, sel ast.SelectionSet, v *models.Logout) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec._Logout(ctx, sel, v)
}

// endregion ***************************** type.gotpl *****************************
