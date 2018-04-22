package shopapi

import (
	"btdxcx.com/os/wrapper"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/micro/jwtauth-srv/wrapper"
	"btdxcx.com/os/custom-error"
	"net/http"
	"github.com/micro/go-log"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"time"
	"github.com/micro/go-web"
	"github.com/micro/cli"
	dproto "btdxcx.com/micro/shop-srv/proto/shop/details"
	kproto "btdxcx.com/micro/shop-srv/proto/shop"
)

const (
	srvShopName = "com.btdxcx.micro.srv.shop"
	apiServiceName = "com.btdxcx.center.api.shops"
)

var (
	detailsCl dproto.ShopClient
	keyCl kproto.ShopKeyClient
)

// Commands add auth api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "shop",
			Usage:  "Run shop api",
			Action: api,
		},
	}
}

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

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "center")
	tokenWrapper := wrapper.NewClientWrapper("center")

	detailsCl = dproto.NewShopClient(
		srvShopName, 
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)
	keyCl = kproto.NewShopKeyClient(srvShopName, client.DefaultClient)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/shops")

	ws.Route(ws.GET("/").To(api.list))
	ws.Route(ws.POST("/").To(api.create))
	ws.Route(ws.GET("/{shopID}").To(api.read))
	ws.Route(ws.PUT("/{shopID}").To(api.update))
	ws.Route(ws.DELETE("/{shopID}").To(api.delete))
	ws.Route(ws.PUT("/{shopID}/submin").To(api.submin))
	ws.Route(ws.PUT("/{shopID}/audit").To(api.audit))
	ws.Route(ws.PUT("/{shopID}/cancel").To(api.cancel))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

// API is APIs
type API struct{}

func (api *API) list(req *restful.Request, rsp *restful.Response) {

	in := new(dproto.ListRequest)
	req.ReadEntity(in)

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	list, err := detailsCl.List(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	
	rsp.WriteEntity(list)
}

func (api *API) create(req *restful.Request, rsp *restful.Response) {

	request := new(dproto.ShopDetails)
	if err1 := req.ReadEntity(request); err1 != nil {
		rsp.WriteError(http.StatusBadRequest, err1)
		return
	}

	in := new(dproto.CreateRequest)
	in.Details = request

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	results, err := detailsCl.Create(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
	rsp.WriteEntity(results)
}

func (api *API) read(req *restful.Request, rsp *restful.Response) {

	in := &dproto.ReadRequest{
		ShopId: req.PathParameter("shopID"),
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	results, err := detailsCl.Read(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteEntity(results)
}

func (api *API) update(req *restful.Request, rsp *restful.Response) {

	request := new(dproto.ShopDetails)
	if err1 := req.ReadEntity(request); err1 != nil {
		rsp.WriteError(http.StatusBadRequest, err1)
		return
	}

	in := new(dproto.UpdateRequest)
	in.ShopId = req.PathParameter("shopID")
	in.Details = request

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	_, err := detailsCl.Update(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

func (api *API) delete(req *restful.Request, rsp *restful.Response) {

	in := &dproto.DeleteRequest{
		ShopId: req.PathParameter("shopID"),
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	_, err := detailsCl.Delete(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

func (api *API) submin(req *restful.Request, rsp *restful.Response) {
	updateState(req, rsp, dproto.State_reviewing)
}

func (api *API) audit(req *restful.Request, rsp *restful.Response) {
	updateState(req, rsp, dproto.State_completed)
}

func (api *API) cancel(req *restful.Request, rsp *restful.Response) {
	updateState(req, rsp, dproto.State_untreated)
}

func updateState(req *restful.Request, rsp *restful.Response, state dproto.State) {
	request := new(dproto.UpdateRequest)
	request.ShopId = req.PathParameter("shopID")
	request.State = state
	
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))
	
	_, err := detailsCl.Update(ctx, request)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}