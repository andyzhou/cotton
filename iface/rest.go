package iface

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
)

//base rest interface
type IRest interface {
	GetPathPara(
		name string,
		req *restful.Request,
	) string
	ParseReqForm(
		formFace interface{},
		req *restful.Request,
	) error
	GetSchemaDecoder() *schema.Decoder
	CreateParameter(
		name, kind, defaultVal string,
		ws *restful.WebService,
	) *restful.Parameter
	CreatePathParameter(
		name, kind string,
		ws *restful.WebService,
	) *restful.Parameter
	RegisterDynamicSubRoute(
		method string,
		consumes string,
		dynamicRootUrl string,
		dynamicPathSlice []string,
		routeFunc restful.RouteFunction,
	) bool
	RegisterSubRoute(
		method, routeUrl, consumes string,
		parameters [] *restful.Parameter,
		routeFunc restful.RouteFunction,
	) bool
}