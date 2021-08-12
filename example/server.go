package main

import (
	"fmt"
	"github.com/andyzhou/cotton"
	"github.com/emicklei/go-restful/v3"
	"os"
	"sync"
)

//server example
const (
	ServerPort = 8080
)

//cb for sub router
func cbOfRouter(req *restful.Request, resp *restful.Response) {
	fmt.Println("cbOfRouter...")
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
