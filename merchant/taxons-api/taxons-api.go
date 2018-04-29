package taxonsapi

import (
	"btdxcx.com/os/wrapper"
	"net/http"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/errors"

	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	proto "btdxcx.com/micro/taxons-srv/proto/taxons"
	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
)

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "taxons",
			Usage:  "Run taxons api",
			Action: api,
		},
	}
}

// API is APIs
type API struct{}

var (
	cl proto.TaxonsClient
)

const (
	serviceName = "com.btdxcx.merchant.api.taxons"
	clientName  = "com.btdxcx.micro.srv.taxons"
)

func (api *API) rootTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContextCustomScope(ctx, "back")

	response, err := cl.RootTaxons(ctx, &proto.RootTaxonsRequest{})
	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response.Message.Children)
}

func (api *API) createTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	record := new(proto.TaxonsRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	request := new(proto.CreateTasonsRequest)
	request.Record = record

	response, err := cl.CreateTaxons(ctx, request)
	if err != nil {
		writeError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(response.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createChildren(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	record := new(proto.TaxonsRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	request := new(proto.CreateChildrenTaxonsRequest)
	request.Code = req.PathParameter("code")
	request.Record = record

	response, err := cl.CreateChildrenTaxons(ctx, request)
	if err != nil {
		writeError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(response.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	record := new(proto.TaxonsRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	request := new(proto.UpdateTaxonsRequest)
	request.Code = req.PathParameter("code")
	request.Record = record

	response, err := cl.UpdateTaxons(ctx, request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(response.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	request := new(proto.DeleteTasonsRequest)
	request.Code = req.PathParameter("code")
	
	response, err := cl.DeleteTaxons(ctx, request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(response); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func writeError(err error, rsp *restful.Response) {
	error := errors.Parse(err.Error())
	if error.Code == 0 {
		rsp.WriteError(500, err)
	} else {
		rsp.WriteError(int(error.Code), error)
	}
}

func api(ctx *cli.Context) {
	service := web.NewService(
		web.Name(serviceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	//	service.Init()

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := jwrapper.NewClientWrapper("back")

	cl = proto.NewTaxonsClient(
		clientName, 
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/taxons")

	ws.Route(ws.GET("").To(api.rootTaxons))
	ws.Route(ws.POST("").To(api.createTaxons))

	ws.Route(ws.POST("/{code}").To(api.createChildren))
	ws.Route(ws.PUT("/{code}").To(api.updateTaxons))
	ws.Route(ws.DELETE("/{code}").To(api.deleteTaxons))

	wc.Add(ws)

	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", proto.TaxonsMessage{})
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}
