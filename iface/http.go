package iface

//http resp interface
type IHttpResp interface {
	GetResp() ([]byte, error)
	SetResp([]byte, error) bool
	SetErr(err error) bool
}

//http req interface
type IHttpReq interface {
	GetResp() (IHttpResp, bool)
	SendResp(IHttpResp) bool
	GetFilePara()(string, string)
	GetBody() []byte
	GetParams() map[string]interface{}
	GetHeaders() map[string]string
	GetReqKind() int
	GetReqUrl() string
	SetFilePara(path, para string) bool
	SetBody(data []byte) bool
	AddParam(key string, val interface{}) bool
	AddHeader(key, val string) bool
	SetReqKind(kind int)
	SetReqUrl(url string) bool
}

//http client queue interface
type IHttpQueue interface {
	Quit()
	SendReq(req IHttpReq) IHttpResp
}
