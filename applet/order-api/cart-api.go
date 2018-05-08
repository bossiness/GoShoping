package orderapi

import (
	"github.com/micro/go-micro/errors"
	"net/http"
	"time"

	wrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/order-srv/proto/order"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/os/custom-error"
	"btdxcx.com/os/wrapper"
	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
)

const (
	cartsAPIServiceName = "com.btdxcx.applet.api.carts"
)

func cart(ctx *cli.Context) {
	service := web.NewService(
		web.Name(cartsAPIServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")
	tokenWrapper := wrapper.NewClientWrapper("mini")

	orderCl = proto.NewOrderClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/carts")

	ws.Route(ws.POST("").To(api.createCart))
	ws.Route(ws.GET("").To(api.readCarts))
	ws.Route(ws.GET("/{id}").To(api.readCart))
	ws.Route(ws.DELETE("/{id}").To(api.deleteCart))

	ws.Route(ws.POST("/{cartId}/items").To(api.createCartItem))
	ws.Route(ws.PATCH("/{cartId}/items/{id}").To(api.modifyCartItem))
	ws.Route(ws.DELETE("/{cartId}/items/{id}").To(api.deleteCartItem))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func (api *API) noop(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	rsp.WriteHeader(http.StatusNotImplemented)
	rsp.WriteEntity("Not Implemented")
}

func (api *API) createCart(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateCartRequest)
	if err := req.ReadEntity(in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}

	out, err := orderCl.CreateCart(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusCreated)
}

func (api *API) readCarts(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadCustomerOrdersRequest)
	in.Customer = ""
	in.State = "cart"

	out, err := orderCl.ReadCustomerOrders(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Records); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusOK)
}

func (api *API) readCart(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadOrderRequest)
	in.Uuid = req.PathParameter("id")

	out, err := orderCl.ReadOrder(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusOK)
}

func (api *API) deleteCart(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteOrderRequest)
	in.Uuid = req.PathParameter("id")

	out, err := orderCl.DeleteOrder(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusNoContent)
}

func (api *API) createCartItem(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateCartItemRequest)
	if err := req.ReadEntity(in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.CartId = req.PathParameter("cartId")
	if len(in.Variant) == 0 {
		rsp.WriteError(http.StatusBadRequest, 
			errors.BadRequest(cartsAPIServiceName + ".createCartItem", "variant is empty."))
		return
	}
	if in.Quantity <= 0 {
		rsp.WriteError(http.StatusBadRequest, 
			errors.BadRequest(cartsAPIServiceName + ".createCartItem", "quantity is fault."))
		return
	}

	out, err := orderCl.CreateCartItem(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Item); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusCreated)
}

func (api *API) modifyCartItem(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateCartItemRequest)
	if err := req.ReadEntity(in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.CartId = req.PathParameter("cartId")
	in.CartItemId = req.PathParameter("id")
	if in.Quantity <= 0 {
		rsp.WriteError(http.StatusBadRequest, 
			errors.BadRequest(cartsAPIServiceName + ".createCartItem", "quantity is fault."))
		return
	}

	out, err := orderCl.UpdateCartItem(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Item); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
	rsp.WriteHeader(http.StatusNoContent)
}

func (api *API) deleteCartItem(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteCartItemRequest)
	if err := req.ReadEntity(in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.CartId = req.PathParameter("cartId")
	in.CartItemId = req.PathParameter("id")

	if _, err := orderCl.DeleteCartItem(ctx, in); err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}