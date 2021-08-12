package face

import (
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
)

/*
 * base rest face
 */

//face info
type Rest struct {
	ws *restful.WebService //reference
	decoder *schema.Decoder
}

//construct
func NewRest(ws *restful.WebService) *Rest {
	//self init
	this := &Rest{
		ws: ws,
		decoder: schema.NewDecoder(),
	}
	this.decoder.IgnoreUnknownKeys(true)
	return this
}

//get schema decoder
func (f *Rest) GetSchemaDecoder() *schema.Decoder {
	return f.decoder
}

//get path param
func (f *Rest) GetPathPara(
				name string,
				req *restful.Request,
			) string {
	return req.PathParameter(name)
}

//parse request form
func (f *Rest) ParseReqForm(
				formFace interface{},
				req *restful.Request,
			) error {
	//basic check
	if formFace == nil || req == nil {
		return errors.New("invalid parameters")
	}

	//parse post form
	err := req.Request.ParseForm()
	if err != nil {
		return err
	}

	//decode form data
	err = f.decoder.Decode(formFace, req.Request.PostForm)
	if err != nil {
		return err
	}
	return nil
}

//register dynamic sub route
//dynamicRootUrl like /test/{para1}/{para2}
//should use request.PathParameter("para1") to get real path parameter value
func (f *Rest) RegisterDynamicSubRoute(
				method string,
				consumes string,
				dynamicRootUrl string,
				dynamicPathSlice []string,
				routeFunc restful.RouteFunction,
			) bool {

	//basic check
	if dynamicRootUrl == "" || routeFunc == nil {
		return false
	}

	//init new route builder
	rb := new(restful.RouteBuilder)
	rb.Method(method).Path(dynamicRootUrl).To(routeFunc)

	if consumes != "" {
		rb.Consumes(consumes)
	}

	//init path parameter
	if dynamicPathSlice != nil && len(dynamicPathSlice) > 0 {
		for _, key := range dynamicPathSlice {
			pp := f.CreatePathParameter(key, "string", f.ws)
			rb.Param(pp)
		}
	}

	//add sub route
	f.ws.Route(rb)
	return true
}

//register static sub route
func (f *Rest) RegisterSubRoute(
				method, routeUrl, consumes string,
				parameters [] *restful.Parameter,
				routeFunc restful.RouteFunction,
			) bool {

	//basic check
	if method == "" || routeUrl == "" || routeFunc == nil {
		return false
	}

	//init new route builder
	rb := new(restful.RouteBuilder)

	//set method, request url and route func
	rb.Method(method).Path(routeUrl).To(routeFunc)

	if consumes != "" {
		rb.Consumes(consumes)
	}

	//set parameter
	if parameters != nil && len(parameters) > 0 {
		for _, parameter := range parameters {
			//set sub parameter
			rb.Param(parameter)
		}
	}

	//add sub route
	f.ws.Route(rb)
	return true
}

//create ws form parameter
func (f *Rest) CreateParameter(
				name, kind, defaultVal string,
				ws *restful.WebService,
			) *restful.Parameter {
	//basic check
	if name == "" || kind == "" {
		return nil
	}
	//init new
	param := ws.FormParameter(name, "").DataType(kind).DefaultValue(defaultVal)
	return param
}

//create ws path parameter
func (f *Rest) CreatePathParameter(
				name, kind string,
				ws *restful.WebService,
			) *restful.Parameter {
	//basic check
	if name == "" || kind == "" {
		return nil
	}
	//init new
	param := ws.PathParameter(name, "").DataType(kind)
	return param
}