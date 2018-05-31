package center

import (
	"btdxcx.com/walker/apis/common/filters"
	"btdxcx.com/walker/apis/common/server"

	"btdxcx.com/walker/apis/common/errors"
	"btdxcx.com/walker/handler"
	"btdxcx.com/walker/model"
	"github.com/emicklei/go-restful"

	"btdxcx.com/walker/repository"
	"btdxcx.com/walker/service/auth"
)

// AuthAPI apis for auth
type AuthAPI struct {
	path    string
	handler *handler.AuthHandler
}

// NewAuthAPI new hello apis
func NewAuthAPI() AuthAPI {
	return AuthAPI{
		path: "/auth",
		handler: &handler.AuthHandler{
			Dot: "center",
		},
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

	api.handler.Service = &auth.Service{
		Session: repository.SingleSession(),
	}

	ws.Route(ws.POST("signin").To(noop).
		// docs
		Doc("登录一个账号").
		Reads(model.AuthRequest{}).
		Writes(model.Token{}).
		Returns(201, "Created", model.Token{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("signup").To(api.handler.Signup).
		// docs
		Doc("用户注册").
		Reads(model.AuthRequest{}).
		Writes(model.Token{}).
		Returns(201, "Created", model.Token{}).
		Returns(401, "Bad Request", nil))

	server.GetContainer().Add(ws)
}

func noop(req *restful.Request, resp *restful.Response) {
	errors.Response(resp, errors.NotImplemented("apis.center.AuthAPI", "not implemented"))
}
