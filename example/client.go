package main

import (
	"fmt"
	"github.com/andyzhou/cotton"
	"github.com/andyzhou/cotton/face"
	"net/http"
	"strings"
)

const (
	reqDomain = "http://localhost:8090"
	reqUrl = "/topic/file"
)

//http client example
func main() {
	//init client
	client := cotton.NewClient()

	//setup request
	req := face.NewHttpReq()
	req.SetReqUrl(fmt.Sprintf("%s%s", reqDomain, reqUrl))
	req.SetReqKind(cotton.ReqKindOfPost)

	//add para
	req.AddParam("jwt", "xxx")

	//send request
	resp := client.SendReq(req)
	body, err := resp.GetResp()
	fmt.Println("body:", string(body))
	fmt.Println("err:", err)
}

func testing()  {
	client := &http.Client{}

	reqUrl := fmt.Sprintf("%s%s", reqDomain, reqUrl)
	fmt.Println(reqUrl)

	req, err := http.NewRequest("POST", reqUrl, strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
