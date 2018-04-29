package shopapi

import (
	"time"

	wrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/shop-srv/proto/shop/details"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/os/custom-error"
	"btdxcx.com/os/wrapper"
	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
)

// Commands add shop api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "shop",
			Usage:  "Run shop api",
			Action: api,
		},
	}
}

const (
	clientName     = "com.btdxcx.micro.srv.shop"
	apiServiceName = "com.btdxcx.applet.api.shops"
)

var (
	shopCl proto.ShopClient
)

func api(ctx *cli.Context) {
	service := web.NewService(
		web.Name(apiServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")
	tokenWrapper := wrapper.NewClientWrapper("mini")

	shopCl = proto.NewShopClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/shop")

	ws.Route(ws.GET("/me").To(api.read))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

// API is APIs
type API struct{}

func (api *API) read(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := &proto.ReadRequest{}
	results, err := shopCl.Read(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteEntity(results)
}
