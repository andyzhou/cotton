package main

import (
	"fmt"
	"github.com/andyzhou/cotton"
	"github.com/andyzhou/cotton/iface"
	"github.com/emicklei/go-restful/v3"
	"io"
	"os"
	"sync"
)

//server example
const (
	ServerPort = 8080
)

//cb for sub router
func cbOfRouter(req *restful.Request, resp *restful.Response, tool iface.ITool) {
	page := req.QueryParameter("page")
	io.WriteString(resp, "cbOfRouter...")
	io.WriteString(resp, fmt.Sprintf("page:%s", page))
}

func main() {
	var (
		wg sync.WaitGroup
	)

	//creat router
	router := cotton.NewRouter(ServerPort)

	//init sub router
	subRouter := &cotton.DynamicSubRoute{
		RouteUrl: "/",
		Module: "test",
		Action: "list",
		RouteFunc: cbOfRouter,
	}
	err := router.RegisterDynamicRoute(subRouter)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	wg.Add(1)
	//start
	router.Start()
	wg.Wait()
}
