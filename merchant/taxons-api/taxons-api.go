package taxonsapi

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-log"

	"time"
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"

	prtot "btdxcx.com/micro/taxons-srv/proto/imp"

)

// API is APIs
type API struct{}

var (
	cl prtot.TaxonsClient
)

const (
	serviceName = "com.btdxcx.merchant.api.taxons"
	clientName = "com.btdxcx.shop.srv.taxons"
)


func (api *API) rootTaxons(req *restful.Request, rsp *restful.Response)  {
	log.Log("Received API.rootTaxons API request")

	shopID := req.HeaderParameter("M-SHOP-ID")
	response, err := cl.Root(context.TODO(), &prtot.TaxonsShopIDRequest{ ShopID: shopID })

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) createTaxons(req *restful.Request, rsp *restful.Response)  {

	shopID := req.HeaderParameter("M-SHOP-ID")
	request := new(prtot.TasonsCreateRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.AddHeader("Content-Type", "text/plain")
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	request.ShopID = shopID

	response, err := cl.Create(context.TODO(), request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) createChildren(req *restful.Request, rsp *restful.Response)  {

	request := new(prtot.TaxonsRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.AddHeader("Content-Type", "text/plain")
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	request.ShopID = req.HeaderParameter("M-SHOP-ID")
	request.Code = req.PathParameter("code")

	response, err := cl.CreateChildren(context.TODO(), request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) updateTaxons(req *restful.Request, rsp *restful.Response)  {

	request := new(prtot.TaxonsRequest)
	if err := req.ReadEntity(&request); err != nil {
		rsp.AddHeader("Content-Type", "text/plain")
		rsp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	request.ShopID = req.HeaderParameter("M-SHOP-ID")
	request.Code = req.PathParameter("code")

	response, err := cl.Update(context.TODO(), request)

	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response)
}

func (api *API) deleteTaxons(req *restful.Request, rsp *restful.Response)  {

	shopID := req.HeaderParameter("M-SHOP-ID")
	code := req.PathParameter("code")
	response, err := cl.Delete(context.TODO(), &prtot.TasonsDeleteRequest{ ShopID: shopID, Code: code })

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

func api(ctx *cli.Context)  {
	service := web.NewService(
		web.Name(serviceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl")) * time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval")) * time.Second,
		),
	)

//	service.Init()

	cl = prtot.NewTaxonsClient(clientName, client.DefaultClient)

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
			Name:  "taxons",
			Usage: "Run taxons api",
			Action: api,
		},
	}
}
