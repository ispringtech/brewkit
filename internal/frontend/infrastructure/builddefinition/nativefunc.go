package builddefinition

import (
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

type nativeFunc interface {
	nativeFunc() *jsonnet.NativeFunction
}

type nativeFunc2[V1, V2 any] struct {
	name string
	v1   argDesc
	v2   argDesc
	f    func(v1 V1, v2 V2) (interface{}, error)
}

func (f nativeFunc2[V1, V2]) nativeFunc() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: f.name,
		Func: errWrapper(
			f.name,
			func(i []interface{}) (interface{}, error) {
				const argsCount = 2
				if len(i) != argsCount {
					return nil, errors.Errorf("not enough arguments to call, expected %d", argsCount)
				}

				v1, err := checkArg[V1](f.v1, i[0])
				if err != nil {
					return nil, err
				}

				v2, err := checkArg[V2](f.v2, i[1])
				if err != nil {
					return nil, err
				}

				return f.f(v1, v2)
			},
		),
		Params: []ast.Identifier{
			ast.Identifier(f.v1.name),
			ast.Identifier(f.v2.name),
		},
	}
}

type nativeFunc3[V1, V2, V3 any] struct {
	name string
	v1   argDesc
	v2   argDesc
	v3   argDesc
	f    func(v1 V1, v2 V2, v3 V3) (interface{}, error)
}

func (f nativeFunc3[V1, V2, V3]) nativeFunc() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: f.name,
		Func: errWrapper(
			f.name,
			func(i []interface{}) (interface{}, error) {
				const argsCount = 3
				if len(i) != argsCount {
					return nil, errors.Errorf("not enough arguments to call, expected %d", argsCount)
				}

				v1, err := checkArg[V1](f.v1, i[0])
				if err != nil {
					return nil, err
				}

				v2, err := checkArg[V2](f.v2, i[1])
				if err != nil {
					return nil, err
				}

				v3, err := checkArg[V3](f.v3, i[2])
				if err != nil {
					return nil, err
				}

				return f.f(v1, v2, v3)
			},
		),
		Params: []ast.Identifier{
			ast.Identifier(f.v1.name),
			ast.Identifier(f.v2.name),
			ast.Identifier(f.v3.name),
		},
	}
}

func errWrapper(name string, next func(i []interface{}) (interface{}, error)) func(i []interface{}) (interface{}, error) {
	return func(i []interface{}) (interface{}, error) {
		json, err := next(i)
		if err != nil {
			err = errors.Wrapf(err, "call '%s' failed", name)
		}
		return json, err
	}
}

type argDesc struct {
	name string
}

// check argument type according to generic type and returns value converted to generic type and error
func checkArg[V any](arg argDesc, argV interface{}) (v V, err error) {
	var ok bool
	v, ok = argV.(V)
	if !ok {
		return v, errors.Errorf("expected '%T' got '%T' as %s arg", v, argV, arg.name)
	}
	return v, nil
}
