package productapi

import (
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
	svrName = "com.btdxcx.applet.api.products-taxon"
)


func taxon(ctx *cli.Context) {
	service := web.NewService(
		web.Name(svrName),
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
	ws.Path("/products-taxon")

	ws.Route(ws.GET("/{taxonCode}").To(api.taxonProducts))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) taxonProducts(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	in := new(proto.TaxonProductsRequest)
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
	in.TaxonCode = req.PathParameter("taxonCode")

	out, err := productCl.TaxonProducts(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

