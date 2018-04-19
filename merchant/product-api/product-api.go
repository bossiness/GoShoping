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
	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
)

const (
	serviceName = "com.btdxcx.merchant.api.product"
	clientName  = "com.btdxcx.micro.srv.product"
)

var (
	attributeCl proto.AttributeClient
	optionCl proto.OptionClient
)

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

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := jwrapper.NewClientWrapper("back")

	attributeCl = proto.NewAttributeClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)
	optionCl = proto.NewOptionClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/product")

	ws.Route(ws.POST("/attributes").To(api.createAttribute))
	ws.Route(ws.GET("/attributes").To(api.readAttributes))
	ws.Route(ws.GET("/attributes/{code}").To(api.readAttribute))
	ws.Route(ws.PUT("/attributes/{code}").To(api.updateAttribute))
	ws.Route(ws.DELETE("/attributes/{code}").To(api.deleteAttribute))

	ws.Route(ws.POST("/options").To(api.createOption))
	ws.Route(ws.GET("/options").To(api.readOptions))
	ws.Route(ws.GET("/options/{code}").To(api.readOption))
	ws.Route(ws.PUT("/options/{code}").To(api.updateOption))
	ws.Route(ws.DELETE("/options/{code}").To(api.deleteOption))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}


func (api *API) createAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateAttributeRequest)
	record := new(proto.AttributesRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := attributeCl.CreateAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readAttributes(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadAttributesRequest)
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

	out, err := attributeCl.ReadAttributes(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadAttributeRequest)
	in.Code = req.PathParameter("code")
	if len(in.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}

	out, err := attributeCl.ReadAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateAttributeRequest)
	record := new(proto.AttributesRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	record.Code = req.PathParameter("code")
	if len(record.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}
	in.Record = record

	out, err := attributeCl.UpdateAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteAttributeRequest)
	in.Code = req.PathParameter("code")
	if len(in.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}

	out, err := attributeCl.DeleteAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createOption(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateOptionRequest)
	record := new(proto.OptionRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := optionCl.CreateOption(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readOptions(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadOptionsRequest)
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

	out, err := optionCl.ReadOptions(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) readOption(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadOptionequest)
	in.Code = req.PathParameter("code")
	if len(in.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}

	out, err := optionCl.ReadOption(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateOption(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateOptionRequest)
	record := new(proto.OptionRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Code  = req.PathParameter("code")
	if len(in.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}
	record.Code = in.Code
	in.Record = record

	out, err := optionCl.UpdateOption(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out.Record); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteOption(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteOptionRequest)
	in.Code = req.PathParameter("code")
	if len(in.Code) == 0 {
		err := errors.BadRequest(serviceName, "code is empty")
		customerror.WriteError(err, rsp)
		return
	}

	out, err := optionCl.DeleteOption(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}
	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}