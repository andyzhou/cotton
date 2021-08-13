package face

import (
	"fmt"
	"github.com/andyzhou/cotton/define"
	"github.com/andyzhou/cotton/iface"
	"github.com/emicklei/go-restful/v3"
	"sync"
)

/*
 * inter router face
 */

//sub router reg info
type SubRouterReg struct {
	Module string
	Action string
	Func func(req *restful.Request, resp *restful.Response)
}

//face info
type Router struct {
	subRouteFace sync.Map //tag -> func
	tool iface.ITool
}

//construct
func NewRouter() *Router {
	//self init
	this := &Router{
		subRouteFace: sync.Map{},
		tool: NewTool(),
	}
	return this
}

//set jwt
func (f *Router) SetJwt(secretKey string) bool {
	return f.tool.SetJwt(secretKey)
}

//entry
func (f *Router) Entry(
					req *restful.Request,
					resp *restful.Response,
					) {
	//get path param
	module := req.PathParameter(define.InterPathParaOfModule)
	action := req.PathParameter(define.InterPathParaOfAction)

	//get cb for sub router
	tag := f.formatTag(module, action)
	v := f.getSubRouteByTag(tag)
	if v == nil {
		return
	}
	//call sub route func
	v(req, resp, f.tool)
}

//register sub router
func (f *Router) RegisterRoute(
				module, action string,
				cb func(*restful.Request, *restful.Response, iface.ITool),
			) bool {
	//check
	if module == "" || cb == nil {
		return false
	}
	tag := f.formatTag(module, action)
	v := f.getSubRouteByTag(tag)
	if v != nil {
		return true
	}
	//create new
	f.subRouteFace.Store(tag, cb)
	return true
}

////////////////
//private func
////////////////

//format tag
func (f *Router) formatTag(module, action string) string {
	return fmt.Sprintf("%s-%s", module, action)
}

//get sub route by tag
func (f *Router) getSubRouteByTag(
					tag string,
				) func(*restful.Request, *restful.Response, iface.ITool) {
	if tag == "" {
		return nil
	}
	v, ok := f.subRouteFace.Load(tag)
	if !ok {
		return nil
	}
	cb, ok := v.(func(*restful.Request, *restful.Response, iface.ITool))
	if !ok {
		return nil
	}
	return cb
}