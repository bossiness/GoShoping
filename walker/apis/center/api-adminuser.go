package center

import (
	"btdxcx.com/walker/apis/common/filters"
	"btdxcx.com/walker/apis/common/server"

	"btdxcx.com/walker/handler"
	"btdxcx.com/walker/model"
	"github.com/emicklei/go-restful"

	"btdxcx.com/walker/repository"
	"btdxcx.com/walker/service/adminuser"
	"btdxcx.com/walker/service/auth"
)

// AdminUserAPI apis for adminuser
type AdminUserAPI struct {
	path    string
	handler *handler.AdminUserHandler
}

// NewAdminUserAPI new hello apis
func NewAdminUserAPI() AdminUserAPI {
	return AdminUserAPI{
		path:    "/users",
		handler: &handler.AdminUserHandler{},
	}
}

// Path get url
func (api *AdminUserAPI) Path() string {
	return api.path
}

// RegisterTo api
func (api AdminUserAPI) RegisterTo(server server.APIServer) {
	ws := new(restful.WebService)

	ws.
		Path(api.Path()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	authService := &auth.Service{
		Session: repository.SingleSession(),
	}

	ws.
		Filter(filter.NCSACommonLogFormatLogger()).
		Filter(filter.CORS(server.GetContainer())).
		Filter(filter.ShopKEYFilter("center")).
		Filter(filter.BearerAuthenticate(authService))

	userService := &adminuser.Service{
		Session: repository.SingleSession(),
	}
	api.handler.Service = userService

	ws.Route(ws.POST("").To(api.handler.Create).
		// docs
		Doc("创建一个管理员用户").
		Reads(model.AdminUser{}).
		Writes(model.NoContent{}).
		Returns(201, "Created", model.NoContent{}))

	ws.Route(ws.GET("").To(api.handler.Reads).
		// docs
		Doc("获取管理员用户集合").
		Returns(200, "Ok", model.AdminUsersPage{}).
		Returns(401, "Bad Request", nil))

	ws.Route(ws.GET("/{id}").To(api.handler.Read).
		// docs
		Doc("获得一个管理员用户").
		Returns(200, "OK", model.AdminUser{}).
		Returns(401, "Bad Request", nil))

	ws.Route(ws.PUT("/{id}").To(api.handler.Update).
		// docs
		Doc("更新一个管理员用户").
		Reads(model.AdminUser{}).
		Returns(204, "No Content", model.AdminUser{}).
		Returns(401, "Bad Request", nil))

	ws.Route(ws.DELETE("/{id}").To(api.handler.Delete).
		// docs
		Doc("删除一个管理员用户").
		Returns(204, "No Content", model.NoContent{}).
		Returns(401, "Bad Request", nil))

	ws.Route(ws.GET("/profile/{username}").To(api.handler.ReadProfile).
		// docs
		Doc("获得一个管理员用户").
		Returns(200, "OK", model.AdminUser{}).
		Returns(401, "Bad Request", nil))

	server.GetContainer().Add(ws)
}
