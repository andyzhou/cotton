package define

const (
	InterRouterName = "router"
	InterPathParaOfModule = "module"
	InterPathParaOfAction = "action"
)

const (
	HttpClientMin = 5
	HttpClientMax = 1024
	HttpClientTimeOut = 10
	HttpReqChanSize = 1024 * 5
)

const (
	HttpReqGet = iota
	HttpReqPost
)