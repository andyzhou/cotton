package main

import (
	"fmt"
	"github.com/andyzhou/cotton"
	"github.com/andyzhou/cotton/face"
)

const (
	reqDomain = "http://localhost:8080"
	reqUrl = "/test/list"
)

//http client example
func main() {
	//init client
	client := cotton.NewClient()

	//setup request
	req := face.NewHttpReq()
	req.SetReqUrl(fmt.Sprintf("%s/%s", reqDomain, reqUrl))

	//send request
	resp := client.SendReq(req)
	body, err := resp.GetResp()
	fmt.Println("body:", string(body))
	fmt.Println("err:", err)
}
