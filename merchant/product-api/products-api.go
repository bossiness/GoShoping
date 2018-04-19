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
	srvName = "com.btdxcx.merchant.api.products"
)

var (
	productCl proto.ProductClient
)

func apis(ctx *cli.Context) {
	service := web.NewService(
		web.Name(srvName),
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
	ws.Route(ws.PUT("/{spu}/associations/{code}").To(api.updateProductAssociation))
	ws.Route(ws.DELETE("/{spu}/associations/{code}").To(api.deleteProductAssociation))

	ws.Route(ws.POST("/{spu}/images").To(api.createProductImage))
	ws.Route(ws.PUT("/{spu}/images/{code}").To(api.updateProductImage))
	ws.Route(ws.DELETE("/{spu}/images/{code}").To(api.deleteProductImage))

	ws.Route(ws.POST("/{spu}/reviews").To(api.createProductReview))
	ws.Route(ws.GET("/{spu}/reviews").To(api.fetchProductReviews))
	ws.Route(ws.PUT("/{spu}/reviews/{id}").To(api.fetchProductReview))
	ws.Route(ws.DELETE("/{spu}/reviews/{id}").To(api.deleteProductReview))
	ws.Route(ws.PATCH("/{spu}/reviews/{id}/accept").To(api.accept))
	ws.Route(ws.PATCH("/{spu}/reviews/{id}/reject").To(api.reject))

	ws.Route(ws.POST("/{spu}/variants").To(api.createProductVariant))
	ws.Route(ws.GET("/{spu}/variants").To(api.fetchProductVariants))
	ws.Route(ws.GET("/{spu}/variants/{sku}").To(api.fetchProductVariant))
	ws.Route(ws.PUT("/{spu}/variants/{sku}").To(api.updateProductVariant))
	ws.Route(ws.DELETE("/{spu}/variants/{sku}").To(api.deleteProductVariant))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) createProduct(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductRequest)
	record := new(proto.ProductRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := productCl.CreateProduct(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProducts(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadProductsRequest)
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
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

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

func (api *API) modifyProduct(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateProductRequest)
	in.Spu = req.PathParameter("spu")
	record := new(proto.ProductRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := productCl.UpdateProduct(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) modifyTaxons(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ModifyTaxonsRequest)
	in.Spu = req.PathParameter("spu")
	taxons := new(proto.ProductTaxon)
	if err := req.ReadEntity(&taxons); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Taxons = taxons

	out, err := productCl.ModifyProductTaxons(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createProductAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductAttributeRequest)
	in.Spu = req.PathParameter("spu")
	attribute := new(proto.ProductAttribute)
	if err := req.ReadEntity(&attribute); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Attribute = attribute

	out, err := productCl.CreateProductAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateProductAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateProductAttributeRequest)
	if err := req.ReadEntity(&in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Spu = req.PathParameter("spu")
	in.Code = req.PathParameter("code")

	out, err := productCl.UpdateProductAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteProductAttribute(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteProductAttributeRequest)
	in.Spu = req.PathParameter("spu")
	in.Code = req.PathParameter("code")

	out, err := productCl.DeleteProductAttribute(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createProductAssociation(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductAssociationRequest)
	in.Spu = req.PathParameter("spu")
	association := new(proto.ProductAssociation)
	if err := req.ReadEntity(&association); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Association = association

	out, err := productCl.CreateProductAssociation(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateProductAssociation(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateProductAssociationRequest)
	in.Spu = req.PathParameter("spu")
	in.Code = req.PathParameter("code")
	association := new(proto.ProductAssociation)
	if err := req.ReadEntity(&association); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Association = association

	out, err := productCl.UpdateProductAssociation(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteProductAssociation(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteProductAssociationRequest)
	in.Spu = req.PathParameter("spu")
	in.Code = req.PathParameter("code")

	out, err := productCl.DeleteProductAssociation(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createProductImage(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductImageRequest)
	in.Spu = req.PathParameter("spu")
	image := new(proto.Image)
	if err := req.ReadEntity(&image); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Image = image

	out, err := productCl.CreateProductImage(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateProductImage(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateProductImageRequest)
	in.Spu = req.PathParameter("spu")
	image := new(proto.Image)
	if err := req.ReadEntity(&image); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Image = image

	out, err := productCl.UpdateProductImage(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteProductImage(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteProductImageRequest)
	in.Spu = req.PathParameter("spu")
	in.Type = req.PathParameter("type")

	out, err := productCl.DeleteProductImage(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createProductReview(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductReviewRequest)
	in.Spu = req.PathParameter("spu")
	record := new(proto.ReviewsRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := productCl.CreateProductReview(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProductReviews(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadProductReviewsRequest)
	in.Spu = req.PathParameter("spu")
	offset, err1 := strconv.Atoi(req.QueryParameter("offset"))
	if err1 != nil {
		rsp.WriteError(
			http.StatusBadRequest, 
			errors.BadRequest(srvName, "query parameter offset error: %s", err1.Error()),
		)
		return
	}
	limit, err2 := strconv.Atoi(req.QueryParameter("limit"))
	if err2 != nil {
		rsp.WriteError(
			http.StatusBadRequest, 
			errors.BadRequest(srvName, "query parameter limit error %s", err2.Error()),
		)
		return
	}
	in.Offset = int32(offset)
	in.Limit = int32(limit)

	out, err := productCl.ReadProductReviews(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProductReview(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadProductReviewRequest)
	in.Spu = req.PathParameter("spu")
	in.Id = req.QueryParameter("id")

	out, err := productCl.ReadProductReview(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteProductReview(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteProductReviewRequest)
	in.Spu = req.PathParameter("spu")
	in.Id = req.QueryParameter("id")

	out, err := productCl.DeleteProductReview(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) accept(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.AcceptProductReviewRequest)
	in.Spu = req.PathParameter("spu")
	in.Id = req.QueryParameter("id")

	out, err := productCl.AcceptProductReview(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) reject(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.RejectProductReviewRequest)
	in.Spu = req.PathParameter("spu")
	in.Id = req.QueryParameter("id")

	out, err := productCl.RejectProductReview(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) createProductVariant(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.CreateProductVariantRequest)
	in.Spu = req.PathParameter("spu")
	record := new(proto.VariantRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := productCl.CreateProductVariant(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProductVariants(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadProductVariantsRequest)
	in.Spu = req.PathParameter("spu")

	out, err := productCl.ReadProductVariants(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) fetchProductVariant(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.ReadProductVariantRequest)
	in.Spu = req.PathParameter("spu")
	in.Sku = req.PathParameter("sku")

	out, err := productCl.ReadProductVariant(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) updateProductVariant(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.UpdateProductVariantRequest)
	in.Spu = req.PathParameter("spu")
	in.Sku = req.PathParameter("sku")
	record := new(proto.VariantRecord)
	if err := req.ReadEntity(&record); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}
	in.Record = record

	out, err := productCl.UpdateProductVariant(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}

func (api *API) deleteProductVariant(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = jwrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	in := new(proto.DeleteProductVariantRequest)
	in.Spu = req.PathParameter("spu")
	in.Sku = req.PathParameter("sku")

	out, err := productCl.DeleteProductVariant(ctx, in)
	if err != nil {
		customerror.WriteError(err, rsp)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	}
}
