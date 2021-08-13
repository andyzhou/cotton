package face

import "github.com/andyzhou/cotton/iface"

/*
 * http request face
 */

//face info
type HttpReq struct {
	Kind int //GET or POST
	Url string
	Headers map[string]string
	Params map[string]interface{}
	Body []byte
	FilePath string
	FilePara string
	ReceiverChan chan iface.IHttpResp `http request receiver chan`
	IsAsync bool
}

//construct
func NewHttpReq() *HttpReq {
	this := &HttpReq{
		Headers:make(map[string]string),
		Params:make(map[string]interface{}),
		Body:make([]byte, 0),
		ReceiverChan:make(chan iface.IHttpResp, 1),
	}
	return this
}

//get sync response
func (f *HttpReq) GetResp() (iface.IHttpResp, bool) {
	resp, ok := <- f.ReceiverChan
	return resp, ok
}

//send sync response
func (f *HttpReq) SendResp(resp iface.IHttpResp) bool {
	if resp == nil || f.ReceiverChan == nil {
		return false
	}
	f.ReceiverChan <- resp
	return true
}

//get file para, return path, para
func (f *HttpReq) GetFilePara()(string, string) {
	return f.FilePath, f.FilePara
}

//set file para
func (f *HttpReq) SetFilePara(path, para string) bool {
	if path == "" || para == "" {
		return false
	}
	f.FilePath = path
	f.FilePara = para
	return true
}

//get body
func (f *HttpReq) GetBody() []byte {
	return f.Body
}

//set body
func (f *HttpReq) SetBody(data []byte) bool {
	if data == nil || len(data) <= 0 {
		return false
	}
	f.Body = data
	return true
}

//get params
func (f *HttpReq) GetParams() map[string]interface{} {
	return f.Params
}

//add param
func (f *HttpReq) AddParam(key string, val interface{}) bool {
	if key == "" || val == nil {
		return false
	}
	f.Params[key] = val
	return true
}

//get headers
func (f *HttpReq) GetHeaders() map[string]string {
	return f.Headers
}

//add header
func (f *HttpReq) AddHeader(key, val string) bool {
	if key == "" || val == "" {
		return false
	}
	f.Headers[key] = val
	return true
}

//get request kind
func (f *HttpReq) GetReqKind() int {
	return f.Kind
}

//set request kind
func (f *HttpReq) SetReqKind(kind int) {
	f.Kind = kind
}

//get request url
func (f *HttpReq) GetReqUrl() string {
	return f.Url
}

//set request url
func (f *HttpReq) SetReqUrl(url string) bool {
	if url == "" {
		return false
	}
	f.Url = url
	return true
}