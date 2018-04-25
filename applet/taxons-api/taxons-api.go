package taxonsapi

import (
	"btdxcx.com/os/wrapper"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/errors"

	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	prtot "btdxcx.com/micro/taxons-srv/proto/imp"
)

// API is APIs
type API struct{}

var (
	cl prtot.TaxonsClient
)

const (
	serviceName = "com.btdxcx.applet.api.taxons"
	clientName  = "com.btdxcx.shop.srv.taxons"
)

func (api *API) rootTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	response, err := cl.Root(ctx, &prtot.TaxonsShopIDRequest{})
	if err != nil {
		writeError(err, rsp)
		return
	}

	rsp.WriteEntity(response.Children)
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

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")

	cl = prtot.NewTaxonsClient(
		clientName, 
		shopkeyWrapper(client.DefaultClient),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/taxons")

	ws.Route(ws.GET("").To(api.rootTaxons))

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
