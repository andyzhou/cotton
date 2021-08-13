package iface

import "github.com/emicklei/go-restful/v3"

//rest tool interface
type ITool interface {
	SetJwt(secretKey string) bool
	GetJwt() IJwt
	ParseReqForm(formFace interface{}, req *restful.Request) error
	GetReferDomain(req *restful.Request) string
	GetClientIp(req *restful.Request) []string
}
