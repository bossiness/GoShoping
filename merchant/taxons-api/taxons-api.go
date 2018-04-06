package taxonsapi

import (
	"net/http"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/errors"

	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	prtot "btdxcx.com/micro/taxons-srv/proto/imp"
	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
)

// API is APIs
type API struct{}

var (
	cl prtot.TaxonsClient
)

const (
	serviceName = "com.btdxcx.merchant.api.taxons"
	clientName  = "com.btdxcx.shop.srv.taxons"
)

func (api *API) rootTaxons(req *restful.Request, rsp *restful.Response) {
	log.Log("Received API.rootTaxons API request")

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	// ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	ctx = jwrapper.NewContextCustomScope(ctx, "back")

	response, err := cl.Root(ctx, &prtot.TaxonsShopIDRequest{})

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response.Children)
}

func (api *API) createTaxons(req *restful.Request, rsp *restful.Response) {

	request := new(prtot.TasonsCreateRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	response, err := cl.Create(ctx, request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) createChildren(req *restful.Request, rsp *restful.Response) {

	request := new(prtot.TaxonsRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	request.Code = req.PathParameter("code")

	response, err := cl.CreateChildren(ctx, request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) updateTaxons(req *restful.Request, rsp *restful.Response) {

	request := new(prtot.TaxonsRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	request.Code = req.PathParameter("code")

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	response, err := cl.Update(ctx, request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) deleteTaxons(req *restful.Request, rsp *restful.Response) {

	code := req.PathParameter("code")

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	response, err := cl.Delete(ctx, &prtot.TasonsDeleteRequest{Code: code})

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
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

	cl = prtot.NewTaxonsClient(
		clientName, 
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
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

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", prtot.TaxonsMessage{})
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}
