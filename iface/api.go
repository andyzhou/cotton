package iface

import "github.com/emicklei/go-restful/v3"

//base api interface
type IApi interface {
	Call(
		name string,
		req *restful.Request,
		resp *restful.Response,
		extPathParams ...interface{},
	) error
	ApiIsExists(name string) bool
	BindApi(name string, face interface{}) error
	SetEntryMethod(name string)
}