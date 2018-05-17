package example

import (
	"io"

	"btdxcx.com/walker/apis/common/filters"
	"btdxcx.com/walker/apis/common/server"

	"github.com/emicklei/go-restful"
)

// HelloAPI apis for hello
type HelloAPI struct {
	path string
}

// NewHelloAPI new hello apis
func NewHelloAPI() HelloAPI {
	return HelloAPI{
		path: "/hello",
	}
}

// Path get url
func (api HelloAPI) Path() string {
	return api.path
}

// RegisterTo api
func (api HelloAPI) RegisterTo(server *server.APIServer) {
	ws := new(restful.WebService)

	ws.
		Path(api.Path()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.
		Filter(filter.NCSACommonLogFormatLogger()).
		Filter(filter.CORS(server.Container()))

	ws.Route(ws.GET("").To(hello).
		// docs
		Doc("hello").
		Writes("world"))

	server.Container().Add(ws)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}
