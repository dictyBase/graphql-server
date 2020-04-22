// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import annotation "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
import context "context"
import empty "github.com/golang/protobuf/ptypes/empty"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"

// TaggedAnnotationServiceClient is an autogenerated mock type for the TaggedAnnotationServiceClient type
type TaggedAnnotationServiceClient struct {
	mock.Mock
}

// AddToAnnotationGroup provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) AddToAnnotationGroup(ctx context.Context, in *annotation.AnnotationGroupId, opts ...grpc.CallOption) (*annotation.TaggedAnnotationGroup, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotationGroup
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.AnnotationGroupId, ...grpc.CallOption) *annotation.TaggedAnnotationGroup); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotationGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.AnnotationGroupId, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAnnotation provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) CreateAnnotation(ctx context.Context, in *annotation.NewTaggedAnnotation, opts ...grpc.CallOption) (*annotation.TaggedAnnotation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotation
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.NewTaggedAnnotation, ...grpc.CallOption) *annotation.TaggedAnnotation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.NewTaggedAnnotation, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAnnotationGroup provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) CreateAnnotationGroup(ctx context.Context, in *annotation.AnnotationIdList, opts ...grpc.CallOption) (*annotation.TaggedAnnotationGroup, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotationGroup
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.AnnotationIdList, ...grpc.CallOption) *annotation.TaggedAnnotationGroup); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotationGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.AnnotationIdList, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAnnotation provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) DeleteAnnotation(ctx context.Context, in *annotation.DeleteAnnotationRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *empty.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.DeleteAnnotationRequest, ...grpc.CallOption) *empty.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*empty.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.DeleteAnnotationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAnnotationGroup provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) DeleteAnnotationGroup(ctx context.Context, in *annotation.GroupEntryId, opts ...grpc.CallOption) (*empty.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *empty.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.GroupEntryId, ...grpc.CallOption) *empty.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*empty.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.GroupEntryId, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAnnotation provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) GetAnnotation(ctx context.Context, in *annotation.AnnotationId, opts ...grpc.CallOption) (*annotation.TaggedAnnotation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotation
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.AnnotationId, ...grpc.CallOption) *annotation.TaggedAnnotation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.AnnotationId, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAnnotationGroup provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) GetAnnotationGroup(ctx context.Context, in *annotation.GroupEntryId, opts ...grpc.CallOption) (*annotation.TaggedAnnotationGroup, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotationGroup
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.GroupEntryId, ...grpc.CallOption) *annotation.TaggedAnnotationGroup); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotationGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.GroupEntryId, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAnnotationTag provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) GetAnnotationTag(ctx context.Context, in *annotation.TagRequest, opts ...grpc.CallOption) (*annotation.AnnotationTag, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.AnnotationTag
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.TagRequest, ...grpc.CallOption) *annotation.AnnotationTag); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.AnnotationTag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.TagRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntryAnnotation provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) GetEntryAnnotation(ctx context.Context, in *annotation.EntryAnnotationRequest, opts ...grpc.CallOption) (*annotation.TaggedAnnotation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotation
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.EntryAnnotationRequest, ...grpc.CallOption) *annotation.TaggedAnnotation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.EntryAnnotationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAnnotationGroups provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) ListAnnotationGroups(ctx context.Context, in *annotation.ListGroupParameters, opts ...grpc.CallOption) (*annotation.TaggedAnnotationGroupCollection, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotationGroupCollection
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.ListGroupParameters, ...grpc.CallOption) *annotation.TaggedAnnotationGroupCollection); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotationGroupCollection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.ListGroupParameters, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAnnotations provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) ListAnnotations(ctx context.Context, in *annotation.ListParameters, opts ...grpc.CallOption) (*annotation.TaggedAnnotationCollection, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotationCollection
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.ListParameters, ...grpc.CallOption) *annotation.TaggedAnnotationCollection); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotationCollection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.ListParameters, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAnnotation provides a mock function with given fields: ctx, in, opts
func (_m *TaggedAnnotationServiceClient) UpdateAnnotation(ctx context.Context, in *annotation.TaggedAnnotationUpdate, opts ...grpc.CallOption) (*annotation.TaggedAnnotation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *annotation.TaggedAnnotation
	if rf, ok := ret.Get(0).(func(context.Context, *annotation.TaggedAnnotationUpdate, ...grpc.CallOption) *annotation.TaggedAnnotation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*annotation.TaggedAnnotation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *annotation.TaggedAnnotationUpdate, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
