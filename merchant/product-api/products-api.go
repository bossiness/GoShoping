package productapi

import (
	"time"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/micro/go-micro/client"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"

	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	srvName = "com.btdxcx.merchant.api.products"
)

var (
	productCl proto.ProductClient
)

func apis(ctx *cli.Context) {
	service := web.NewService(
		web.Name(serviceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := jwrapper.NewClientWrapper("back")

	productCl = proto.NewProductClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/products")

	ws.Route(ws.POST("").To(api.createProduct))
	ws.Route(ws.GET("").To(api.fetchProducts))
	ws.Route(ws.GET("/{spu}").To(api.fetchProduct))
	ws.Route(ws.PATCH("/{spu}").To(api.modifyProduct))
	ws.Route(ws.PATCH("/{spu}/taxons").To(api.modifyTaxons))

	ws.Route(ws.POST("/{spu}/attributes").To(api.createProductAttribute))
	ws.Route(ws.PUT("/{spu}/attributes/{code}").To(api.updateProductAttribute))
	ws.Route(ws.DELETE("/{spu}/attributes/{code}").To(api.deleteProductAttribute))

	ws.Route(ws.POST("/{spu}/associations").To(api.createProductAssociation))
	ws.Route(ws.PUT("/{spu}/associations/{code}").To(api.noop))
	ws.Route(ws.DELETE("/{spu}/associations/{code}").To(api.noop))

	ws.Route(ws.POST("/{spu}/images").To(api.noop))
	ws.Route(ws.PUT("/{spu}/images/{code}").To(api.noop))
	ws.Route(ws.DELETE("/{spu}/images/{code}").To(api.noop))

	ws.Route(ws.POST("/{spu}/reviews").To(api.noop))
	ws.Route(ws.GET("/{spu}/reviews").To(api.noop))
	ws.Route(ws.PUT("/{spu}/reviews/{id}").To(api.noop))
	ws.Route(ws.DELETE("/{spu}/reviews/{id}").To(api.noop))
	ws.Route(ws.PATCH("/{spu}/reviews/{id}/accept").To(api.noop))
	ws.Route(ws.PATCH("/{spu}/reviews/{id}/reject").To(api.noop))

	ws.Route(ws.POST("/{spu}/variants").To(api.noop))
	ws.Route(ws.GET("/{spu}/variants").To(api.noop))
	ws.Route(ws.GET("/{spu}/variants/{sku}").To(api.noop))
	ws.Route(ws.PUT("/{spu}/variants/{sku}").To(api.noop))
	ws.Route(ws.DELETE("/{spu}/variants/{sku}").To(api.noop))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) noop(req *restful.Request, rsp *restful.Response) {
}

