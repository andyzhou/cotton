package cotton

import (
	"errors"
	"fmt"
	"github.com/andyzhou/cotton/define"
	"github.com/andyzhou/cotton/face"
	"github.com/andyzhou/cotton/iface"
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

//macro define
const (
	//req method
	HttpReqGet = "GET"
	HttpReqPost = "POST"

	//req form type
	HttpReqTypeTextPlain = "text/plain"
	HttpReqTypeFormEncode = "application/x-www-form-urlencoded"
)

//dynamic sub route info
type DynamicSubRoute struct {
	RouteUrl string
	Module string
	Action string
	ReqMethod string //default 'GET'
	ReqType string //default 'text/plain'
	Consumes string //optional
	PathSlice []string //optional
	RouteFunc func(
					request *restful.Request,
					response *restful.Response,
				)
}

//router info
type Router struct {
	httpPort int
	ws *restful.WebService
	rest iface.IRest
	router iface.IRouter
}

//construct, step-1
func NewRouter(httpPort int) *Router {
	//init ws
	ws := new(restful.WebService)
	restful.Add(ws)

	//self init
	this := &Router{
		httpPort: httpPort,
		ws: ws,
		rest: face.NewRest(ws),
		router: face.NewRouter(),
	}
	this.interInit()
	return this
}

//register dynamic route, step-2
func (r *Router) RegisterDynamicRoute(router *DynamicSubRoute) error {
	//check
	if router == nil || router.RouteUrl == "" || router.RouteFunc == nil {
		return errors.New("invalid parameter")
	}
	if router.ReqMethod == "" {
		router.ReqMethod = HttpReqGet
	}

	//setup sub router request url
	dynamicReqUrl := fmt.Sprintf("%s/{%s}/{%s}",
							router.RouteUrl,
							define.InterPathParaOfModule,
							define.InterPathParaOfAction,
						)
	//init dynamic path slice
	dynamicPathSlice := []string{
		define.InterPathParaOfModule,
		define.InterPathParaOfAction,
	}

	//register sub route
	r.rest.RegisterDynamicSubRoute(
				router.ReqMethod,
				router.Consumes,
				dynamicReqUrl,
				dynamicPathSlice,
				r.router.Entry,
			)

	//register sub router
	r.router.RegisterRoute(
				router.Module,
				router.Action,
				router.RouteFunc,
			)

	return nil
}

//start service, step-3
func (r *Router) Start() {
	addr := fmt.Sprintf(":%d", r.httpPort)
	go http.ListenAndServe(addr, nil)
}

//parse req form
func (r *Router) ParseReqForm(
				formFace interface{},
				req *restful.Request,
			) error {
	return r.rest.ParseReqForm(formFace, req)
}

///////////////
//private func
///////////////

//inter init
func (r *Router) interInit() {
}