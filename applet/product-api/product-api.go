package productapi

import (
	"github.com/micro/go-micro/errors"
	"strconv"
	"btdxcx.com/os/custom-error"
	"net/http"
	"github.com/micro/go-log"
	"github.com/emicklei/go-restful"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/client"
	"time"
	"github.com/micro/go-web"
	"github.com/micro/cli"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	serviceName = "com.btdxcx.applet.api.products"
	clientName  = "com.btdxcx.micro.srv.product"
)

var (
	productCl proto.ProductClient
)

// API is APIs
type API struct{}

// Commands add command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "product",
			Usage:  "Run product api",
			Action: api,
		},
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

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")

	productCl = proto.NewProductClient(
		clientName,
		shopkeyWrapper(client.DefaultClient),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/products")

	ws.Route(ws.GET("").To(api.fetchProducts))
	ws.Route(ws.GET("/{spu}").To(api.fetchProduct))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) fetchProducts(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	in := new(proto.ReadProductsRequest)
	offset, err1 := strconv.Atoi(req.QueryParameter("offset"))
	if err1 != nil {
		rsp.WriteError(
			http.StatusBadRequest, 
			errors.BadRequest(serviceName, "query parameter limit error: %s", err1.Error()),
		)
		return
	}
	limit, err2 := strconv.Atoi(req.QueryParameter("limit"))
	if err2 != nil {
		rsp.WriteError(
			http.StatusBadRequest, 
			errors.BadRequest(serviceName, "query parameter limit error %s", err2.Error()),
		)
		return
	}
	in.Offset = int32(offset)
	in.Limit = int32(limit)

	out, err := productCl.ReadProducts(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProduct(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	in := new(proto.ReadProductRequest)
	in.Spu = req.PathParameter("spu")

	out, err := productCl.ReadProduct(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}
