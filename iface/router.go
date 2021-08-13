package iface

import (
	"github.com/emicklei/go-restful/v3"
)

//router interface
type IRouter interface {
	Entry(
		req *restful.Request,
		resp *restful.Response,
	)
	RegisterRoute(
		module, action string,
		cb func(req *restful.Request, resp *restful.Response, tool ITool),
	) bool
	SetJwt(secretKey string) bool
}
