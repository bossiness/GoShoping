package center

import (
	"io"

	"btdxcx.com/walker/apis/common/filters"
	"btdxcx.com/walker/apis/common/server"

	"github.com/emicklei/go-restful"
	"btdxcx.com/walker/apis/common/errors"
	"btdxcx.com/walker/model"
)

// AuthAPI apis for auth
type AuthAPI struct {
	path string
}

// NewAuthAPI new hello apis
func NewAuthAPI() AuthAPI {
	return AuthAPI{
		path: "/auth",
	}
}

// Path get url
func (api *AuthAPI) Path() string {
	return api.path
}

// RegisterTo api
func (api AuthAPI) RegisterTo(server server.APIServer) {
	ws := new(restful.WebService)

	ws.
		Path(api.Path()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.
		Filter(filter.NCSACommonLogFormatLogger()).
		Filter(filter.CORS(server.GetContainer()))

	ws.Route(ws.POST("signin").To(noop).
		// docs
		Doc("登录一个账号").
		Reads(model.AuthRequest{}).
		Writes(model.Token{}).
		Returns(201, "Created", model.Token{}).
		Returns(404, "Not Found", nil))

	server.GetContainer().Add(ws)
}

func noop(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
	errors.NotImplemented("apis.center.AuthAPI", "not implemented")
}
