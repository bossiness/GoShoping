package memberapi

import (
	"net/http"
	"strconv"
	"time"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/os/custom-error"
	"btdxcx.com/os/wrapper"
	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/member-srv/proto/member"
)

const (
	usersAPIService = "com.btdxcx.merchant.api.users"
)

var (
	adminUserCl proto.AdminUserClient
)

func adminuserAPIs(ctx *cli.Context) {
	service := web.NewService(
		web.Name(usersAPIService),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := jwrapper.NewClientWrapper("back")

	customerCl = proto.NewCustomerClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/users")

	ws.Route(ws.POST("").To(api.createAdminUser))
	ws.Route(ws.GET("").To(api.readAdminUsers))
	ws.Route(ws.GET("/{id}").To(api.readAdminUser))
	ws.Route(ws.PUT("/{id}").To(api.updateAdminUser))
	ws.Route(ws.DELETE("/{id}").To(api.deleteAdminUser))
	ws.Route(ws.GET("/profile/{name}").To(api.readAdminUserFromName))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func (api *API) createAdminUser(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateAdminUserRequest)
	record := new(proto.AdminUserRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := adminUserCl.CreateAdminUser(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readAdminUsers(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadAdminUsersRequest)
	offset, err1 := strconv.Atoi(req.QueryParameter("offset"))
	if err1 != nil {
		offset = 0
	}
	limit, err2 := strconv.Atoi(req.QueryParameter("limit"))
	if err2 != nil {
		limit = 20
	}
	in.Offset = int32(offset)
	in.Limit = int32(limit)

	out, err := adminUserCl.ReadAdminUsers(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readAdminUser(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadAdminUserRequest)
	in.Id = req.PathParameter("id")

	out, err := adminUserCl.ReadAdminUser(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readAdminUserFromName(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadAdminUserFormNameRequest)
	in.Name = req.PathParameter("name")

	out, err := adminUserCl.ReadAdminUserFormName(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateAdminUser(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateAdminUserRequest)
	in.Id = req.PathParameter("id")
	record := new(proto.AdminUserRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := adminUserCl.UpdateAdminUser(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteAdminUser(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteAdminUserRequest)
	in.Id = req.PathParameter("id")

	out, err := adminUserCl.DeleteAdminUser(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}
