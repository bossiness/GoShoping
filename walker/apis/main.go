package main

import (
	"btdxcx.com/walker/apis/common/examples"
	"btdxcx.com/walker/apis/common/server"
	"btdxcx.com/walker/apis/common/static"
)


func main() {

	apis := server.NewAPIServer()

	static.NewServeStatic().RegisterTo(apis)
	example.NewHelloAPI().RegisterTo(apis)

	apis.Start()

}