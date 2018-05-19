package main

import (
	"os"
	"btdxcx.com/walker/apis/common/examples"
	"btdxcx.com/walker/apis/common/static"
	"btdxcx.com/walker/apis/common/server"

	"btdxcx.com/walker/apis/applet"
	"btdxcx.com/walker/apis/center"
	"github.com/micro/go-log"
)

func main() {

	if len(os.Args) != 2 {
		log.Log("apis [applet|center|test]")
		return
	}

	if os.Args[1] == "applet" {
		apis := applet.NewAPIServer()

		static.NewServeStatic().RegisterTo(apis)
		example.NewHelloAPI().RegisterTo(apis)

		apis.Start()
	} else if os.Args[1] == "center" {
		apis := center.NewAPIServer()

		center.NewAuthAPI().RegisterTo(apis)

		apis.Start()
	} else if os.Args[1] == "test" {
		apis := server.NewCommonAPIServer()

		static.NewServeStatic().RegisterTo(apis)
		example.NewHelloAPI().RegisterTo(apis)

		apis.Start()
	}

	log.Log("apis [applet|test]")
}
