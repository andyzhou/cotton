package face

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/andyzhou/cotton/define"
	"github.com/andyzhou/cotton/iface"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
 * http face
 */

//face info
type HttpQueue struct {
	clientNum int
	clients map[int]*http.Client //http client instance map
	reqChan chan iface.IHttpReq //request lazy chan
	closeChan chan bool
	sync.RWMutex
}

//construct
func NewHttpQueue(clientNum ...int) *HttpQueue {
	var (
		clientNumVal int
	)

	//set http client num
	clientNumVal = define.HttpClientMin
	if clientNum != nil && len(clientNum) > 0 {
		clientNumVal = clientNum[0]
	}
	if clientNumVal > define.HttpClientMax {
		clientNumVal = define.HttpClientMax
	}

	//self init
	this := &HttpQueue{
		clientNum: clientNumVal,
		clients: make(map[int]*http.Client),
		reqChan: make(chan iface.IHttpReq, define.HttpReqChanSize),
		closeChan: make(chan bool, 1),
	}
	this.interInit()
	go this.runMainProcess()
	return this
}

//quit
func (f *HttpQueue) Quit() {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	f.closeChan <- true
}

//send request
func (f *HttpQueue) SendReq(req iface.IHttpReq) (resp iface.IHttpResp) {
	var (
		isOk bool
	)
	//init resp
	resp = NewHttpResp()
	if req == nil {
		resp.SetErr(errors.New("invalid parameter"))
		return
	}
	defer func() {
		if e := recover(); e != nil {
			resp.SetErr(fmt.Errorf("%v", e))
			return
		}
	}()

	//send to chan
	f.reqChan <- req

	//wait for resp
	resp, isOk = req.GetResp()
	if !isOk {
		resp.SetErr(errors.New("no any response"))
	}
	return
}

/////////////////
//private func
/////////////////

//run main process
func (f *HttpQueue) runMainProcess() {
	var (
		req iface.IHttpReq
		isOk bool
	)

	//defer
	defer func() {
		if e := recover(); e != nil {
			log.Println("HttpQueue:mainProcess panic, err:", e)
		}
		close(f.reqChan)
		close(f.closeChan)
	}()

	//loop
	for {
		select {
		case req, isOk = <- f.reqChan:
			if isOk {
				f.sendRealHttpReq(req)
			}
		case <- f.closeChan:
			return
		}
	}
}

//send real http request
func (f *HttpQueue) sendRealHttpReq(req iface.IHttpReq) bool {
	var (
		httpReq *http.Request
		orgResp *http.Response
		respBody []byte
		err error
	)

	//init resp
	httpResp := NewHttpResp()

	//check
	if req == nil {
		return false
	}

	//get upload file path, para
	filePath, _ := req.GetFilePara()
	if filePath != "" {
		//file upload request
		httpReq, err = f.fileUploadRequest(req)
	}else{
		//general http request
		httpReq, err = f.genHttpRequest(req)
	}
	if err != nil {
		httpResp.Err = err
		req.SendResp(httpResp)
		return false
	}

	//set headers
	headers := req.GetHeaders()
	if headers != nil {
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
	}

	//set http connect close
	httpReq.Header.Set("Connection", "close")
	httpReq.Close = true

	//get http client
	httpClient := f.getClient()
	if httpClient == nil {
		httpResp.Err = errors.New("can't get http client")
		req.SendResp(httpResp)
		return false
	}

	//begin send request
	orgResp, err = httpClient.Do(httpReq)
	if err != nil {
		httpResp.Err = err
		req.SendResp(httpResp)
		return false
	}

	//close resp before return
	defer orgResp.Body.Close()

	//read response
	respBody, err = ioutil.ReadAll(orgResp.Body)
	if err != nil {
		httpResp.Err = err
		req.SendResp(httpResp)
		return false
	}

	//send to response chan
	httpResp.Data = respBody
	req.SendResp(httpResp)
	return true
}

//send file upload request
func (f *HttpQueue) fileUploadRequest(req iface.IHttpReq) (*http.Request, error) {
	var (
		httpRequest *http.Request
		filePart io.Writer
	)

	//get upload path, para
	filePath, filePara := req.GetFilePara()

	//try open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//init multi part
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, err = writer.CreateFormFile(filePara,filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(filePart, file)

	//add extend parameters
	params := req.GetParams()
	for key, val := range params {
		v2, ok := val.(string)
		if !ok {
			continue
		}
		_ = writer.WriteField(key, v2)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	//get req url
	reqUrl := req.GetReqUrl()

	//init http request
	httpRequest, err = http.NewRequest("POST", reqUrl, body)
	httpRequest.Header.Set("Content-Type", writer.FormDataContentType())
	return httpRequest, err
}

//gen http request
func (f *HttpQueue) genHttpRequest(req iface.IHttpReq) (*http.Request, error) {
	var (
		tempStr string
		buffer = bytes.NewBuffer(nil)
		httpReq *http.Request
		err error
	)

	//setup params
	params := req.GetParams()
	if len(params) > 0 {
		i := 0
		for k, v := range params {
			if i > 0 {
				buffer.WriteString("&")
			}
			tempStr = fmt.Sprintf("%s=%v", k, v)
			buffer.WriteString(tempStr)
			i++
		}
	}

	//get request url
	reqUrl := req.GetReqUrl()

	//request kind
	kind := req.GetReqKind()
	switch kind {
	case define.HttpReqPost:
		{
			//check body
			body := req.GetBody()
			if body != nil {
				buffer.Write(body)
			}

			//int post req
			httpReq, err = http.NewRequest("POST", reqUrl, strings.NewReader(buffer.String()))
			if err != nil {
				return nil, err
			}

			//headers
			headers := req.GetHeaders()
			if headers == nil || len(headers) <= 0 {
				httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			//format post form
			if params != nil {
				httpReq.Form = make(url.Values)
				for k, v := range params {
					keyVal := fmt.Sprintf("%v", v)
					httpReq.Form.Add(k, keyVal)
				}
			}
		}
	default:
		//init get req
		if buffer != nil && buffer.Len() > 0 {
			reqUrl = fmt.Sprintf("%s?%s", reqUrl, buffer.String())
		}
		httpReq, err = http.NewRequest("GET", reqUrl, nil)
	}

	return httpReq, err
}

//get rand http client
func (f *HttpQueue) getClient() *http.Client {
	if f.clients == nil {
		return nil
	}
	idx := rand.Intn(f.clientNum) + 1
	if idx > define.HttpClientMax {
		idx = define.HttpClientMax
	}
	v, ok := f.clients[idx]
	if !ok {
		return nil
	}
	return v
}

//create one http client
func (f *HttpQueue) createClient() *http.Client {
	//init http trans
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: define.HttpClientTimeOut * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: define.HttpClientTimeOut * time.Second,
		ResponseHeaderTimeout: define.HttpClientTimeOut * time.Second,
		ExpectContinueTimeout: define.HttpClientTimeOut* time.Second,
	}

	//init native http client
	client := &http.Client{
		Timeout:time.Second * define.HttpClientTimeOut,
		Transport:netTransport,
	}
	return client
}

//inter init
func (f *HttpQueue) interInit() {
	//init batch http clients
	for i := 1; i <= f.clientNum; i++ {
		client := f.createClient()
		f.clients[i] = client
	}
}
