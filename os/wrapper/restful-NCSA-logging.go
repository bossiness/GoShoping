package logwrapper

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"
	"strings"
	"time"
)

// This example shows how to create a filter that produces log lines
// according to the Common Log Format, also known as the NCSA standard.
//
// kindly contributed by leehambley
//
// GET http://localhost:8080/ping

func NCSACommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)
		logstr := fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %d %d",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
		)
		log.Log(logstr)
	}
}
