package static

import (
	"fmt"
	"net/http"
	"path"

	"btdxcx.com/walker/apis/common/server"
	"btdxcx.com/walker/apis/common/static/image"

	"github.com/emicklei/go-restful"
)

// This example shows how to define methods that serve static files
// It uses the standard http.ServeFile method
//
// GET http://localhost:8080/static/resources/test.xml
// GET http://localhost:8080/static/resources/
//
// GET http://localhost:8080/static/resources?resource=subdir/test.xml

var rootdir = "/tmp"

// ServeStatic apis for hello
type ServeStatic struct {
	path string
}

// NewServeStatic new static serve
func NewServeStatic() ServeStatic {
	return ServeStatic{
		path: "/static",
	}
}

// Path get url
func (s ServeStatic) Path() string {
	return s.path
}

// RegisterTo static serve
func (s ServeStatic) RegisterTo(server server.APIServer) {
	ws := new(restful.WebService)
	ws.Path(s.Path())
	
	ws.Route(ws.GET("/resources/{subpath:*}").To(staticFromPathParam))
	ws.Route(ws.GET("/resources").To(staticFromQueryParam))

	ws.Route( ws.GET("/image").To(image.Home))
	ws.Route(ws.POST("/image").To(image.Upload))
	ws.Route( ws.GET("/images/{imgid}").To(image.Download))
	ws.Route( ws.GET("/images").To(image.List))

	server.GetContainer().Add(ws)

	println("[go-restful] serving files on http://localhost:3001/static from local /tmp")
}

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	actual := path.Join(rootdir, req.PathParameter("subpath"))
	fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}
