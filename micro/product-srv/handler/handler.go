package handler

import (
	"btdxcx.com/micro/product-srv/db"
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	svrName = "btdxcx.com/micro/product-srv"
)

// Handler product handler
type Handler struct{}

// CreateProduct is a single request handler called via client.CreateProduct or the generated client code
func (h *Handler) CreateProduct(ctx context.Context, req *proto.CreateProductRequest, rsp *proto.CreateProductResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.CreateProduct(shopID, req.Record)
	if err1 != nil {
		return errors.InternalServerError(svrName + ".CreateProduct", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadProducts is a single request handler called via client.ReadProducts or the generated client code
func (h *Handler) ReadProducts(ctx context.Context, req *proto.ReadProductsRequest, rsp *proto.ReadProductsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	products, err := db.ReadProducts(shopID, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound(svrName + ".ReadProducts", err.Error())
	}

	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*products))
	rsp.Records = *products
	return nil
}

// ReadProduct is a single request handler called via client.ReadProduct or the generated client code
func (h *Handler) ReadProduct(ctx context.Context, req *proto.ReadProductRequest, rsp *proto.ReadProductResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	product, err := db.ReadProduct(shopID, req.Spu)
	if err != nil {
		return errors.NotFound(svrName + ".ReadProduct", err.Error())
	}
	rsp.Record = product
	return nil
}

// UpdateProduct is a single request handler called via client.UpdateProduct or the generated client code
func (h *Handler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest, rsp *proto.UpdateProductResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProduct(shopID, req.Spu, req.Record.Name, req.Record.Description); err != nil {
		return errors.NotFound(svrName + ".UpdateProduct", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteProduct is a single request handler called via client.DeleteProduct or the generated client code
func (h *Handler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest, rsp *proto.DeleteProductResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProduct(shopID, req.Spu); err != nil {
		return errors.NotFound(svrName + ".DeleteProduct", err.Error())
	}

	return nil
}

// ModifyProductTaxons is a single request handler called via client.ModifyProductTaxons or the generated client code
func (h *Handler) ModifyProductTaxons(ctx context.Context, req *proto.ModifyTaxonsRequest, rsp *proto.ModifyTaxonsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProductTaxons(shopID, req.Spu, req.Taxons.Main, req.Taxons.Others); err != nil {
		return errors.NotFound(svrName + ".ModifyProductTaxons", err.Error())
	}

	rsp.Taxons = req.Taxons
	return nil
}

// CreateProductAttribute is a single request handler called via client.CreateProductAttribute or the generated client code
func (h *Handler) CreateProductAttribute(ctx context.Context, req *proto.CreateProductAttributeRequest, rsp *proto.CreateProductAttributeResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.CreateProductAttribute(shopID, req.Spu, req.Attribute); err != nil {
		return errors.NotFound(svrName + ".CreateProductAttribute", err.Error())
	}

	rsp.Attribute = req.Attribute
	return nil
}

// ReadProductAttributes is a single request handler called via client.ReadProductAttributes or the generated client code
func (h *Handler) ReadProductAttributes(ctx context.Context, req *proto.ReadProductAttributesRequest, rsp *proto.ReadProductAttributesResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	attributes, err := db.ReadProductAttributes(shopID, req.Spu)
	if err != nil {
		return errors.NotFound(svrName + ".ReadProductAttributes", err.Error())
	}

	rsp.Attributes = *attributes
	return nil
}

// UpdateProductAttribute is a single request handler called via client.UpdateProductAttribute or the generated client code
func (h *Handler) UpdateProductAttribute(ctx context.Context, req *proto.UpdateProductAttributeRequest, rsp *proto.UpdateProductAttributeResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	record := proto.ProductAttribute{Code: req.Code, Value: req.Value}
	if err := db.UpdateProductAttribute(shopID, req.Spu, &record); err != nil {
		return errors.NotFound(svrName + ".UpdateProductAttribute", err.Error())
	}

	rsp.Attribute = &record
	return nil
}

// DeleteProductAttribute is a single request handler called via client.DeleteProductAttribute or the generated client code
func (h *Handler) DeleteProductAttribute(ctx context.Context, req *proto.DeleteProductAttributeRequest, rsp *proto.DeleteProductAttributeResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductAttribute(shopID, req.Spu, req.Code); err != nil {
		return errors.NotFound(svrName + ".DeleteProductAttribute", err.Error())
	}

	return nil
}

// CreateProductAssociation is a single request handler called via client.CreateProductAssociation or the generated client code
func (h *Handler) CreateProductAssociation(ctx context.Context, req *proto.CreateProductAssociationRequest, rsp *proto.CreateProductAssociationResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}
	
	if err := db.CreateProductAssociation(shopID, req.Spu, req.Association); err != nil {
		return errors.NotFound(svrName + ".CreateProductAssociation", err.Error())
	}

	rsp.Association = req.Association
	return nil
}

func (h *Handler) ReadProductAssociations(ctx context.Context, req *proto.ReadProductAssociationsRequest, rsp *proto.ReadProductAssociationsResponse) error {
	// shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	// if err1 != nil {
	// 	return err1
	// }

	return nil
}

// UpdateProductAssociation is a single request handler called via client.UpdateProductAssociation or the generated client code
func (h *Handler) UpdateProductAssociation(ctx context.Context, req *proto.UpdateProductAssociationRequest, rsp *proto.UpdateProductAssociationResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProductAssociation(shopID, req.Spu, req.Association); err != nil {
		return errors.NotFound(svrName + ".UpdateProductAssociation", err.Error())
	}

	rsp.Association = req.Association
	return nil
}

// DeleteProductAssociation is a single request handler called via client.DeleteProductAssociation or the generated client code
func (h *Handler) DeleteProductAssociation(ctx context.Context, req *proto.DeleteProductAssociationRequest, rsp *proto.DeleteProductAssociationResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductAssociation(shopID, req.Spu, req.Code); err != nil {
		return errors.NotFound(svrName + ".DeleteProductAssociation", err.Error())
	}

	return nil
}

// CreateProductImage is a single request handler called via client.CreateProductImage or the generated client code
func (h *Handler) CreateProductImage(ctx context.Context, req *proto.CreateProductImageRequest, rsp *proto.CreateProductImageResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}
	
	if err := db.CreateProductImage(shopID, req.Spu, req.Image); err != nil {
		return errors.NotFound(svrName + ".CreateProductImage", err.Error())
	}

	rsp.Image = req.Image
	return nil
}
func (h *Handler) ReadProductImages(ctx context.Context, req *proto.ReadProductImagesRequest, rsp *proto.ReadProductImagesResponse) error {
	// shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	// if err1 != nil {
	// 	return err1
	// }


	return nil
}

// UpdateProductImage is a single request handler called via client.UpdateProductImage or the generated client code
func (h *Handler) UpdateProductImage(ctx context.Context, req *proto.UpdateProductImageRequest, rsp *proto.UpdateProductImageResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}
	
	if err := db.UpdateProductImage(shopID, req.Spu, req.Image); err != nil {
		return errors.NotFound(svrName + ".UpdateProductImage", err.Error())
	}

	rsp.Image = req.Image
	return nil
}

// DeleteProductImage is a single request handler called via client.DeleteProductImage or the generated client code
func (h *Handler) DeleteProductImage(ctx context.Context, req *proto.DeleteProductImageRequest, rsp *proto.DeleteProductImageResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductImage(shopID, req.Spu, req.Type); err != nil {
		return errors.NotFound(svrName + ".DeleteProductImage", err.Error())
	}

	return nil
}

// CreateProductReview is a single request handler called via client.CreateProductReview or the generated client code
func (h *Handler) CreateProductReview(ctx context.Context, req *proto.CreateProductReviewRequest, rsp *proto.CreateProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.CreateProductReview(shopID, req.Spu, req.Record); err != nil {
		return errors.NotFound(svrName + ".CreateProductReview", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadProductReviews is a single request handler called via client.ReadProductReviews or the generated client code
func (h *Handler) ReadProductReviews(ctx context.Context, req *proto.ReadProductReviewsRequest, rsp *proto.ReadProductReviewsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	reviews, err := db.ReadProductReviews(shopID, req.Spu, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound(svrName + ".ReadProductReviews", err.Error())
	}

	rsp.Records = *reviews
	rsp.Offset = req.Offset
	rsp.Limit = req.Limit
	rsp.Total = int32(len(*reviews))
	return nil
}

// ReadProductReview is a single request handler called via client.ReadProductReview or the generated client code
func (h *Handler) ReadProductReview(ctx context.Context, req *proto.ReadProductReviewRequest, rsp *proto.ReadProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	review, err := db.ReadProductReview(shopID, req.Id)
	if err != nil {
		return errors.NotFound(svrName + ".ReadProductReview", err.Error())
	}

	rsp.Record = review
	return nil
}

// UpdateProductReview is a single request handler called via client.UpdateProductReview or the generated client code
func (h *Handler) UpdateProductReview(ctx context.Context, req *proto.UpdateProductReviewRequest, rsp *proto.UpdateProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProductReview(shopID, req.Record); err != nil {
		return errors.NotFound(svrName + ".UpdateProductReview", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteProductReview is a single request handler called via client.DeleteProductReview or the generated client code
func (h *Handler) DeleteProductReview(ctx context.Context, req *proto.DeleteProductReviewRequest, rsp *proto.DeleteProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductReview(shopID, req.Id); err != nil {
		return errors.NotFound(svrName + ".DeleteProductReview", err.Error())
	}

	return nil
}

// AcceptProductReview is a single request handler called via client.AcceptProductReview or the generated client code
func (h *Handler) AcceptProductReview(ctx context.Context, req *proto.AcceptProductReviewRequest, rsp *proto.AcceptProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.AcceptProductReview(shopID, req.Id); err != nil {
		return errors.NotFound(svrName + ".AcceptProductReview", err.Error())
	}
	return nil
}

// RejectProductReview is a single request handler called via client.RejectProductReview or the generated client code
func (h *Handler) RejectProductReview(ctx context.Context, req *proto.RejectProductReviewRequest, rsp *proto.RejectProductReviewResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.RejectProductReview(shopID, req.Spu ,req.Id); err != nil {
		return errors.NotFound(svrName + ".RejectProductReview", err.Error())
	}

	return nil
}

// CreateProductOption is a single request handler called via client.CreateProductOption or the generated client code
func (h *Handler) CreateProductOption(ctx context.Context, req *proto.CreateProductOptionRequest, rsp *proto.CreateProductOptionesponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.CreateProductOption(shopID, req.Spu, req.Option); err != nil {
		return errors.NotFound(svrName + ".CreateProductOption", err.Error())
	}

	return nil
}

// ReadProductOptions is a single request handler called via client.ReadProductOptions or the generated client code
func (h *Handler) ReadProductOptions(ctx context.Context, req *proto.ReadProductOptionsRequest, rsp *proto.ReadProductOptionsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	options, err := db.ReadProductOptions(shopID, req.Spu)
	if err != nil {
		return errors.NotFound(svrName + ".DeleteProductOption", err.Error())
	}

	rsp.Options = *options
	return nil
}

// UpdateProductOption is a single request handler called via client.UpdateProductOption or the generated client code
func (h *Handler) UpdateProductOption(ctx context.Context, req *proto.UpdateProductOptionRequest, rsp *proto.UpdateProductOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProductOption(shopID, req.Spu, req.Option); err != nil {
		return errors.NotFound(svrName + ".UpdateProductOption", err.Error())
	}
	return nil
}

// DeleteProductOption is a single request handler called via client.DeleteProductOption or the generated client code
func (h *Handler) DeleteProductOption(ctx context.Context, req *proto.DeleteProductOptionRequest, rsp *proto.DeleteProductOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductOption(shopID, req.Spu, req.Code); err != nil {
		return errors.NotFound(svrName + ".DeleteProductOption", err.Error())
	}
	return nil
}

// CreateProductVariant is a single request handler called via client.CreateProductVariant or the generated client code
func (h *Handler) CreateProductVariant(ctx context.Context, req *proto.CreateProductVariantRequest, rsp *proto.CreateProductVariantReponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.CreateProductVariant(shopID, req.Spu, req.Record); err != nil {
		return errors.NotFound(svrName + ".CreateProductVariant", err.Error())
	}

	return nil
}

// ReadProductVariants is a single request handler called via client.ReadProductVariants or the generated client code
func (h *Handler) ReadProductVariants(ctx context.Context, req *proto.ReadProductVariantsRequest, rsp *proto.ReadProductVariantsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	variants, err := db.ReadProductVariants(shopID, req.Spu, 0, 20)
	if err != nil {
		return errors.NotFound(svrName + ".ReadProductVariants", err.Error())
	}

	rsp.Records = *variants
	return nil
}

// ReadProductVariant is a single request handler called via client.ReadProductVariant or the generated client code
func (h *Handler) ReadProductVariant(ctx context.Context, req *proto.ReadProductVariantRequest, rsp *proto.ReadProductVariantResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	variatnt, err := db.ReadProductVariant(shopID, req.Spu, req.Sku)
	if err != nil {
		return errors.NotFound(svrName + ".ReadProductVariant", err.Error())
	}

	rsp.Record = variatnt
	return nil
}

// UpdateProductVariant is a single request handler called via client.UpdateProductVariant or the generated client code
func (h *Handler) UpdateProductVariant(ctx context.Context, req *proto.UpdateProductVariantRequest, rsp *proto.UpdateProductVariantResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateProductVariant(shopID, req.Spu, req.Sku, req.Record); err != nil {
		return errors.NotFound(svrName + ".UpdateProductVariant", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteProductVariant is a single request handler called via client.DeleteProductVariant or the generated client code
func (h *Handler) DeleteProductVariant(ctx context.Context, req *proto.DeleteProductVariantRequest, rsp *proto.DeleteProductVariantResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteProductVariant(shopID, req.Spu, req.Sku); err != nil {
		return errors.NotFound(svrName + ".DeleteProductVariant", err.Error())
	}

	return nil
}
