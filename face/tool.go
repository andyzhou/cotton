package face

import (
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
	"strings"
)

/*
 * rest tool face
 */

//inter macro define
const (
	HttpProtocol = "://"
)

//face info
type Tool struct {
	decoder *schema.Decoder
}

//construct
func NewTool() *Tool {
	//self init
	this := &Tool{
		decoder: schema.NewDecoder(),
	}
	this.decoder.IgnoreUnknownKeys(true)
	return this
}

//parse request form
func (f *Tool) ParseReqForm(formFace interface{}, req *restful.Request) error {
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

//get refer domain
func (f *Tool) GetReferDomain(referUrl string) string {
	var (
		referDomain string
	)
	if referUrl == "" {
		return referDomain
	}
	//find first '://' pos
	protocolLen := len(HttpProtocol)
	protocolPos := strings.Index(referUrl, HttpProtocol)
	if protocolPos <= -1 {
		return referDomain
	}
	//pick domain
	tempBytes := []byte(referUrl)
	tempBytesLen := len(tempBytes)
	prefixLen := protocolPos + protocolLen
	resetUrl := tempBytes[prefixLen:tempBytesLen]
	tempSlice := strings.Split(string(resetUrl), "/")
	if tempSlice == nil || len(tempSlice) <= 0 {
		return referDomain
	}
	referDomain = fmt.Sprintf("%s%s", tempBytes[0:prefixLen], tempSlice[0])
	return referDomain
}

//get all ip from client
func (f *Tool) GetClientIp(req *restful.Request) []string {
	var (
		tempStr string
		ipSlice = make([]string, 0)
	)

	//get original data
	clientAddress := req.Request.RemoteAddr
	xRealIp := req.Request.Header.Get("X-Real-IP")
	xForwardedFor := req.Request.Header.Get("X-Forwarded-For")

	//analyze general ip
	if clientAddress != "" {
		tempStr = f.analyzeClientIp(clientAddress)
		if tempStr != "" {
			ipSlice = append(ipSlice, tempStr)
		}
	}

	//analyze x-real-ip
	if xRealIp != "" {
		tempStr = f.analyzeClientIp(clientAddress)
		if tempStr != "" {
			ipSlice = append(ipSlice, tempStr)
		}
	}

	//analyze x-forward-for
	//like:192.168.0.1,192.168.0.2
	if xForwardedFor != "" {
		tempSlice := strings.Split(xForwardedFor, ",")
		if len(tempSlice) > 0 {
			for _, tmpAddr := range tempSlice {
				tempStr = f.analyzeClientIp(tmpAddr)
				if tempStr != "" {
					ipSlice = append(ipSlice, tempStr)
				}
			}
		}
	}

	return ipSlice
}

//analyze client ip
func (f *Tool) analyzeClientIp(address string) string {
	tempSlice := strings.Split(address, ":")
	if len(tempSlice) < 1 {
		return ""
	}
	return tempSlice[0]
}