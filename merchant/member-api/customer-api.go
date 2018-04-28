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
	cutomersAPIService = "com.btdxcx.merchant.api.cutomers"
	clientName         = "com.btdxcx.micro.srv.member"
)

var (
	customerCl proto.CustomerClient
)

func cutomersAPIs(ctx *cli.Context) {
	service := web.NewService(
		web.Name(cutomersAPIService),
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
	ws.Path("/cutomers")

	ws.Route(ws.POST("").To(api.createCustomer))
	ws.Route(ws.GET("").To(api.readCustomers))
	ws.Route(ws.GET("/{id}").To(api.readCustomer))
	ws.Route(ws.PUT("/{id}").To(api.updateCustomer))
	ws.Route(ws.DELETE("/{code}").To(api.deleteCustomer))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func (api *API) createCustomer(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateCustomerRequest)
	record := new(proto.CustomerRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := customerCl.CreateCustomer(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readCustomers(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadCustomersRequest)
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

	out, err := customerCl.ReadCustomers(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readCustomer(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadCustomerRequest)
	in.Id = req.PathParameter("id")

	out, err := customerCl.ReadCustomer(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateCustomer(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateCustomerRequest)
	in.Id = req.PathParameter("id")
	record := new(proto.CustomerRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := customerCl.UpdateCustomer(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteCustomer(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteCustomerRequest)
	in.Id = req.PathParameter("id")

	out, err := customerCl.DeleteCustomer(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}
