package iface

import "github.com/emicklei/go-restful/v3"

//rest tool interface
type ITool interface {
	GetJson() IJson
	GetJwt() IJwt
	SetJwt(secretKey string) bool
	ParseReqForm(formFace interface{}, req *restful.Request) error
	GetReferDomain(req *restful.Request) string
	GetClientIp(req *restful.Request) []string
}
