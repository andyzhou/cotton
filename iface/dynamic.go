package iface

import "reflect"

//dynamic faces interface
type IDynamic interface {
	Cast(
			method string,
			params ...interface{},
		) bool
	Call(
			name, method string,
			params ...interface{},
		) ([]reflect.Value, error)
	GetFace(name string) interface{}
	BindFace(name string, face interface{}) bool
}