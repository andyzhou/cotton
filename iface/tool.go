package iface

import "github.com/emicklei/go-restful/v3"

//rest tool interface
type ITool interface {
	ParseReqForm(formFace interface{}, req *restful.Request) error
	GetReferDomain(referUrl string) string
	GetClientIp(req *restful.Request) []string
}
