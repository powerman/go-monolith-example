// Package reflectx provide helpers for reflect.
package reflectx

import (
	"reflect"
	"runtime"
	"strings"
)

// MethodsOf require pointer to interface (e.g.: new(app.Appl)) and
// returns all it methods.
func MethodsOf(v interface{}) []string {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		panic("require pointer to interface")
	}
	typ = typ.Elem()
	methods := make([]string, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		methods[i] = typ.Method(i).Name
	}
	return methods
}

// CallerMethodName returns caller's method name for given stack depth.
func CallerMethodName(skip int) string {
	pc, _, _, _ := runtime.Caller(1 + skip)
	names := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return names[len(names)-1]
}

// CallerPkg returns caller's package name (from path) for given stack depth.
func CallerPkg(skip int) string {
	pc, _, _, _ := runtime.Caller(1 + skip)
	names := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pkg := names[len(names)-1]
	return pkg[:strings.Index(pkg, ".")]
}
