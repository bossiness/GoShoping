package mongodb

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	productsCollectionName = "products"
	reviewsCollectionName  = "reviews"
	variantsCollectionName = "variants"
)

type Product struct {
	ID            bson.ObjectId         `bson:"_id,omitempty"`
	Spu           string                `bson:"spu,omitempty"`
	Name          string                `bson:"name,omitempty"`
	Description   string                `bson:"description,omitempty"`
	MainTaxon     string                `bson:"mainTaxon,omitempty"`
	ProductTaxons []string              `bson:"productTaxons,omitempty"`
	UpdatedAt     int64                 `bson:"updated_at,omitempty"`
	CreatedAt     int64                 `bson:"created_at,omitempty"`
	Associations  []Product_Association `bson:"associations,omitempty"`
	Images        []Product_Image       `bson:"images,omitempty"`
	Attributes    []Product_Attribute   `bson:"attributes,omitempty"`
	Options       []Product_Option      `bson:"options,omitempty"`
}

type Product_Association struct {
	Code string   `bson:"code,omitempty"`
	Spus []string `bson:"spus,omitempty"`
}

type Product_Image struct {
	Type  string   `bson:"type,omitempty"`
	Paths []string `bson:"paths,omitempty"`
}

type Product_Attribute struct {
	Code  string `bson:"code,omitempty"`
	Value string `bson:"value,omitempty"`
}

type Product_Option struct {
	Code  string `bson:"code,omitempty"`
	Value string `bson:"value,omitempty"`
}

type Variant struct {
	ID               bson.ObjectId   `bson:"id,omitempty"`
	Sku              string          `bson:"sku,omitempty"`
	Spu              string          `bson:"spu,omitempty"`
	Name             string          `bson:"name,omitempty"`
	Pricings         Variant_Pricing `bson:"pricings,omitempty"`
	Tracked          bool            `bson:"tracked,omitempty"`
	ShippingCategory string          `bson:"shippingCategory,omitempty"`
	OptionValues     []string        `bson:"optionValues,omitempty"`
	Stock            int32           `bson:"stock,omitempty"`
	Images           []Variant_Image `bson:"images,omitempty"`
	UpdatedAt        int64           `bson:"updated_at,omitempty"`
	CreatedAt        int64           `bson:"created_at,omitempty"`
}

type Variant_Pricing struct {
	Current int32 `bson:"current,omitempty"`
}

type Variant_Image struct {
	Type  string   `bson:"type,omitempty"`
	Paths []string `bson:"paths,omitempty"`
}

type Review struct {
	ID        bson.ObjectId  `bson:"id,omitempty"`
	Spu       string         `bson:"spu,omitempty"`
	Title     string         `bson:"title,omitempty"`
	Comment   string         `bson:"comment,omitempty"`
	Rating    int32          `bson:"rating,omitempty"`
	Author    Reviews_Author `bson:"author,omitempty"`
	Status    string         `bson:"status,omitempty"`
	UpdatedAt int64          `bson:"updated_at,omitempty"`
	CreatedAt int64          `bson:"created_at,omitempty"`
}

type Reviews_Author struct {
	Username string `bson:"username,omitempty"`
	Nick     string `bson:"nick,omitempty"`
}

// CreateProduct Insert
func (m *Mongo) CreateProduct(dbname string, record *proto.ProductRecord) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	index := mgo.Index{
		Key:        []string{"spu"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	associations := []Product_Association{}
	for _, it := range record.Associations {
		item := Product_Association{
			Code: it.Code,
			Spus: it.Spus,
		}
		associations = append(associations, item)
	}

	images := []Product_Image{}
	for _, it := range record.Images {
		item := Product_Image{
			Type:  it.Type,
			Paths: it.Paths,
		}
		images = append(images, item)
	}

	attributes := []Product_Attribute{}
	for _, it := range record.Attributes {
		item := Product_Attribute{
			Code:  it.Code,
			Value: it.Value,
		}
		attributes = append(attributes, item)
	}

	options := []Product_Option{}
	for _, it := range record.Options {
		item := Product_Option{
			Code:  it.Code,
			Value: it.Value,
		}
		options = append(options, item)
	}

	vc := m.session.DB(dbname).C(variantsCollectionName)
	for _, it := range record.Variants {
		pricing := Variant_Pricing{
			Current: it.Pricings.Current,
		}
		vimages := []Variant_Image{}
		for _, i := range it.Images {
			item := Variant_Image{
				Type:  i.Type,
				Paths: i.Paths,
			}
			vimages = append(vimages, item)
		}
		variant := Variant{
			ID:               bson.NewObjectId(),
			Sku:              it.Sku,
			Spu:              record.Spu,
			Name:             it.Name,
			Pricings:         pricing,
			Tracked:          it.Tracked,
			ShippingCategory: it.ShippingCategory,
			OptionValues:     it.OptionValues,
			Stock:            it.Stock,
			Images:           vimages,
			CreatedAt:        time.Now().Unix(),
			UpdatedAt:        time.Now().Unix(),
		}
		vc.Insert(variant)
	}

	rc := m.session.DB(dbname).C(reviewsCollectionName)
	for _, it := range record.Reviews {
		author := Reviews_Author{
			Username: it.Author.Username,
			Nick:     it.Author.Nick,
		}
		review := Review{
			ID:        bson.NewObjectId(),
			Spu:       record.Spu,
			Title:     it.Title,
			Comment:   it.Comment,
			Rating:    it.Rating,
			Author:    author,
			Status:    "new",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		rc.Insert(review)
	}

	doc := Product{
		ID:            bson.NewObjectId(),
		Spu:           record.Spu,
		Name:          record.Name,
		Description:   record.Description,
		MainTaxon:     record.MainTaxon,
		ProductTaxons: record.ProductTaxons,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		Associations:  associations,
		Images:        images,
		Attributes:    attributes,
		Options:       options,
	}

	return c.Insert(&doc)
}

// ReadProducts Find
func (m *Mongo) ReadProducts(dbname string, offset int, limit int) (*[]*proto.ProductRecord, error) {
	c := m.session.DB(dbname).C(productsCollectionName)

	results := []Product{}
	if err := c.Find(nil).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.ProductRecord{}
	for _, it := range results {
		records = append(records, m.readProduct(&it, dbname))
	}

	return &records, nil
}

// ReadProduct Find
func (m *Mongo) ReadProduct(dbname string, spu string) (*proto.ProductRecord, error) {
	c := m.session.DB(dbname).C(productsCollectionName)

	result := Product{}
	selector := bson.M{"spu": spu}
	if err := c.Find(selector).One(&result); err != nil {
		return nil, err
	}

	return m.readProduct(&result, dbname), nil
}

func (m *Mongo) readProduct(it *Product, dbname string) *proto.ProductRecord {
	associations := []*proto.ProductRecord_Association{}
	for _, ia := range it.Associations {
		association := &proto.ProductRecord_Association{
			Code: ia.Code,
			Spus: ia.Spus,
		}
		associations = append(associations, association)
	}

	images := []*proto.Image{}
	for _, ii := range it.Images {
		image := &proto.Image{
			Type:  ii.Type,
			Paths: ii.Paths,
		}
		images = append(images, image)
	}

	attributes := []*proto.ProductRecord_AttributesRecord{}
	for _, ia := range it.Attributes {
		if a, err := m.ReadAttribute(dbname, ia.Code); err == nil {
			attribute := &proto.ProductRecord_AttributesRecord{
				Code:          ia.Code,
				Value:         ia.Value,
				Type:          a.Type,
				Configuration: a.Configuration,
			}
			attributes = append(attributes, attribute)
		}
	}

	options := []*proto.ProductRecord_OptionRecord{}
	for _, io := range it.Options {
		if op, err := m.ReadOption(dbname, io.Code); err == nil {
			poos := []*proto.ProductRecord_OptionRecord_OptionValue{}
			for _, o := range op.Options {
				poo := &proto.ProductRecord_OptionRecord_OptionValue{
					Value:       o.Value,
					Description: o.Description,
				}
				poos = append(poos, poo)
			}
			option := &proto.ProductRecord_OptionRecord{
				Code:    io.Code,
				Value:   io.Code,
				Options: poos,
			}
			options = append(options, option)
		}
	}

	variants := []*proto.VariantRecord{}
	vc := m.session.DB(dbname).C(variantsCollectionName)
	vs := []Variant{}
	if err := vc.Find(bson.M{"spu": it.Spu}).All(&vs); err == nil {
		for _, iv := range vs {
			pricings := &proto.VariantRecord_Pricing{
				Current: iv.Pricings.Current,
			}
			images := []*proto.Image{}
			for _, ii := range iv.Images {
				image := &proto.Image{
					Type:  ii.Type,
					Paths: ii.Paths,
				}
				images = append(images, image)
			}
			variant := &proto.VariantRecord{
				Sku:              iv.Sku,
				Name:             iv.Name,
				Pricings:         pricings,
				Tracked:          iv.Tracked,
				ShippingCategory: iv.ShippingCategory,
				OptionValues:     iv.OptionValues,
				Stock:            iv.Stock,
				Images:           images,
				UpdatedAt:        iv.UpdatedAt,
				CreatedAt:        iv.CreatedAt,
			}
			variants = append(variants, variant)
		}
	}

	reviews := []*proto.ReviewsRecord{}
	rc := m.session.DB(dbname).C(reviewsCollectionName)
	rs := []Review{}
	if err := rc.Find(bson.M{"spu": it.Spu}).All(&rs); err == nil {
		for _, ir := range rs {
			author := &proto.ReviewsRecord_Author{
				Username: ir.Author.Username,
				Nick:     ir.Author.Nick,
			}
			review := &proto.ReviewsRecord{
				Id:        ir.ID.Hex(),
				Title:     ir.Title,
				Comment:   ir.Comment,
				Rating:    ir.Rating,
				Author:    author,
				Status:    ir.Status,
				UpdatedAt: ir.UpdatedAt,
				CreatedAt: ir.CreatedAt,
			}
			reviews = append(reviews, review)
		}
	}

	item := &proto.ProductRecord{
		Spu:           it.Spu,
		Name:          it.Name,
		Description:   it.Description,
		MainTaxon:     it.MainTaxon,
		ProductTaxons: it.ProductTaxons,
		UpdatedAt:     it.UpdatedAt,
		CreatedAt:     it.CreatedAt,
		Associations:  associations,
		Images:        images,
		Attributes:    attributes,
		Options:       options,
		Variants:      variants,
		Reviews:       reviews,
	}
	return item
}

// UpdateProduct Update
func (m *Mongo) UpdateProduct(dbname string, spu string, name string, description string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	updataData := bson.M{"$set": bson.M{
		"name":        name,
		"description": description,
		"updated_at":  time.Now()}}

	return c.Update(selector, updataData)
}

// DeleteProduct Remove
func (m *Mongo) DeleteProduct(dbname string, spu string) error {
	c := m.session.DB(dbname).C(productsCollectionName)
	selector := bson.M{"spu": spu}
	return c.Remove(selector)
}

// UpdateProductTaxons Update
func (m *Mongo) UpdateProductTaxons(dbname string, spu string, main string, others []string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	updataData := bson.M{"$set": bson.M{
		"mainTaxon":     main,
		"productTaxons": others,
		"updated_at":    time.Now()}}

	return c.Update(selector, updataData)
}

// CreateProductAttribute Update
func (m *Mongo) CreateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	attribute := Product_Attribute{
		Code:  record.Code,
		Value: record.Value,
	}
	updataData := bson.M{"$push": bson.M{
		"attributes": attribute}}

	return c.Update(selector, updataData)
}

// ReadProductAttributes Find
func (m *Mongo) ReadProductAttributes(dbname string, spu string) (*[]*proto.ProductRecord_AttributesRecord, error) {

	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return nil, err
	}

	attributes := []*proto.ProductRecord_AttributesRecord{}
	for _, ia := range product.Attributes {
		if a, err := m.ReadAttribute(dbname, ia.Code); err == nil {
			attribute := &proto.ProductRecord_AttributesRecord{
				Code:          ia.Code,
				Value:         ia.Value,
				Type:          a.Type,
				Configuration: a.Configuration,
			}
			attributes = append(attributes, attribute)
		}
	}
	
	return &attributes, nil
}

// UpdateProductAttribute Update
func (m *Mongo) UpdateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Attributes {
		if it.Code == record.Code {
			pullData := bson.M{"$pull": bson.M{
				"attributes": it}}
			if err := c.Update(selector, pullData); err != nil {
				return err
			}

			attribute := Product_Attribute{
				Code:  record.Code,
				Value: record.Value,
			}
			pushData := bson.M{"$push": bson.M{
				"attributes": attribute}}

			return c.Update(selector, pushData)
		}
	}

	return errors.New("没有找到 Code:" + record.Code)
}

// DeleteProductAttribute Update
func (m *Mongo) DeleteProductAttribute(dbname string, spu string, code string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Attributes {
		if it.Code == code {
			pullData := bson.M{"$pull": bson.M{
				"attributes": it}}
			return c.Update(selector, pullData)
		}
	}

	return errors.New("没有找到 Code:" + code)
}

// CreateProductAssociation Update
func (m *Mongo) CreateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	association := Product_Association{
		Code: record.Code,
		Spus: record.Spus,
	}
	updataData := bson.M{"$push": bson.M{
		"associations": association}}

	return c.Update(selector, updataData)
}

// UpdateProductAssociation Update
func (m *Mongo) UpdateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Associations {
		if it.Code == record.Code {
			pullData := bson.M{"$pull": bson.M{
				"associations": it}}
			if err := c.Update(selector, pullData); err != nil {
				return err
			}

			association := Product_Association{
				Code: record.Code,
				Spus: record.Spus,
			}
			pushData := bson.M{"$push": bson.M{
				"associations": association}}

			return c.Update(selector, pushData)
		}
	}

	return errors.New("没有找到 Code:" + record.Code)
}

// DeleteProductAssociation Update
func (m *Mongo) DeleteProductAssociation(dbname string, spu string, code string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Associations {
		if it.Code == code {
			pullData := bson.M{"$pull": bson.M{
				"associations": it}}
			return c.Update(selector, pullData)
		}
	}

	return errors.New("没有找到 Code:" + code)
}

// CreateProductImage Update
func (m *Mongo) CreateProductImage(dbname string, spu string, record *proto.Image) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	image := Product_Image{
		Type:  record.Type,
		Paths: record.Paths,
	}
	updataData := bson.M{"$push": bson.M{
		"images": image}}

	return c.Update(selector, updataData)
}

// UpdateProductImage Update
func (m *Mongo) UpdateProductImage(dbname string, spu string, record *proto.Image) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Images {
		if it.Type == record.Type {
			pullData := bson.M{"$pull": bson.M{
				"images": it}}
			if err := c.Update(selector, pullData); err != nil {
				return err
			}

			image := Product_Image{
				Type:  record.Type,
				Paths: record.Paths,
			}
			pushData := bson.M{"$push": bson.M{
				"images": image}}

			return c.Update(selector, pushData)
		}
	}

	return errors.New("没有找到 Type:" + record.Type)
}

// DeleteProductImage Update
func (m *Mongo) DeleteProductImage(dbname string, spu string, itype string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Images {
		if it.Type == itype {
			pullData := bson.M{"$pull": bson.M{
				"images": it}}
			return c.Update(selector, pullData)
		}
	}

	return errors.New("没有找到 Type:" + itype)
}

// CreateProductReview Update
func (m *Mongo) CreateProductReview(dbname string, spu string, record *proto.ReviewsRecord) error {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	author := Reviews_Author{
		Username: record.Author.Username,
		Nick:     record.Author.Nick,
	}
	review := Review{
		ID:        bson.NewObjectId(),
		Spu:       spu,
		Title:     record.Title,
		Comment:   record.Comment,
		Rating:    record.Rating,
		Author:    author,
		Status:    "new",
		UpdatedAt: time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
	}

	return rc.Insert(review)
}

// ReadProductReviews Find
func (m *Mongo) ReadProductReviews(dbname string, spu string, offset int, limit int) (*[]*proto.ReviewsRecord, error) {
	rc := m.session.DB(dbname).C(reviewsCollectionName)
	selector := bson.M{"spu": spu}

	reviews := []Review{}
	if err := rc.Find(selector).Skip(offset).Limit(limit).All(&reviews); err != nil {
		return nil, err
	}

	records := []*proto.ReviewsRecord{}
	for _, ir := range reviews {
		author := &proto.ReviewsRecord_Author{
			Username: ir.Author.Username,
			Nick:     ir.Author.Nick,
		}
		review := &proto.ReviewsRecord{
			Id:        ir.ID.Hex(),
			Title:     ir.Title,
			Comment:   ir.Comment,
			Rating:    ir.Rating,
			Author:    author,
			Status:    ir.Status,
			UpdatedAt: ir.UpdatedAt,
			CreatedAt: ir.CreatedAt,
		}
		records = append(records, review)
	}

	return &records, nil
}

// ReadProductReview Find
func (m *Mongo) ReadProductReview(dbname string, id string) (*proto.ReviewsRecord, error) {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	review := Review{}
	if err := rc.FindId(bson.ObjectIdHex(id)).One(&review); err != nil {
		return nil, err
	}

	author := &proto.ReviewsRecord_Author{
		Username: review.Author.Username,
		Nick:     review.Author.Nick,
	}
	record := &proto.ReviewsRecord{
		Id:        review.ID.Hex(),
		Title:     review.Title,
		Comment:   review.Comment,
		Rating:    review.Rating,
		Author:    author,
		Status:    review.Status,
		UpdatedAt: review.UpdatedAt,
		CreatedAt: review.CreatedAt,
	}

	return record, nil
}

// UpdateProductReview Update
func (m *Mongo) UpdateProductReview(dbname string, record *proto.ReviewsRecord) error {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	update := bson.M{"$set": bson.M{
		"title":           record.Title,
		"comment":         record.Comment,
		"rating":          record.Rating,
		"author.username": record.Author.Username,
		"author.nick":     record.Author.Nick,
		"updated_at":      time.Now().Unix(),
	}}
	return rc.UpdateId(bson.ObjectIdHex(record.Id), update)
}

// DeleteProductReview Remove id
func (m *Mongo) DeleteProductReview(dbname string, id string) error {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	return rc.RemoveId(bson.ObjectIdHex(id))
}

// AcceptProductReview Update
func (m *Mongo) AcceptProductReview(dbname string, id string) error {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	update := bson.M{"$set": bson.M{
		"status":     "accept",
		"updated_at": time.Now().Unix(),
	}}
	return rc.UpdateId(bson.ObjectIdHex(id), update)
}

// RejectProductReview Update
func (m *Mongo) RejectProductReview(dbname string, spu string, id string) error {
	rc := m.session.DB(dbname).C(reviewsCollectionName)

	update := bson.M{"$set": bson.M{
		"status":     "reject",
		"updated_at": time.Now().Unix(),
	}}
	return rc.UpdateId(bson.ObjectIdHex(id), update)
}

// CreateProductVariant Insert
func (m *Mongo) CreateProductVariant(dbname string, spu string, record *proto.VariantRecord) error {
	vc := m.session.DB(dbname).C(variantsCollectionName)

	pricing := Variant_Pricing{
		Current: record.Pricings.Current,
	}
	vimages := []Variant_Image{}
	for _, i := range record.Images {
		item := Variant_Image{
			Type:  i.Type,
			Paths: i.Paths,
		}
		vimages = append(vimages, item)
	}
	variant := Variant{
		ID:               bson.NewObjectId(),
		Sku:              record.Sku,
		Spu:              spu,
		Name:             record.Name,
		Pricings:         pricing,
		Tracked:          record.Tracked,
		ShippingCategory: record.ShippingCategory,
		OptionValues:     record.OptionValues,
		Stock:            record.Stock,
		Images:           vimages,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
	}
	return vc.Insert(variant)
}

// ReadProductVariants Find
func (m *Mongo) ReadProductVariants(dbname string, spu string, offset int, limit int) (*[]*proto.VariantRecord, error) {
	vc := m.session.DB(dbname).C(variantsCollectionName)

	variants := []Variant{}
	if err := vc.Find(bson.M{"spu": spu}).Skip(offset).Limit(limit).All(&variants); err != nil {
		return nil, err
	}
	records := []*proto.VariantRecord{}
	for _, iv := range variants {
		pricings := &proto.VariantRecord_Pricing{
			Current: iv.Pricings.Current,
		}
		images := []*proto.Image{}
		for _, ii := range iv.Images {
			image := &proto.Image{
				Type:  ii.Type,
				Paths: ii.Paths,
			}
			images = append(images, image)
		}
		variant := &proto.VariantRecord{
			Sku:              iv.Sku,
			Name:             iv.Name,
			Pricings:         pricings,
			Tracked:          iv.Tracked,
			ShippingCategory: iv.ShippingCategory,
			OptionValues:     iv.OptionValues,
			Stock:            iv.Stock,
			Images:           images,
			UpdatedAt:        iv.UpdatedAt,
			CreatedAt:        iv.CreatedAt,
		}
		records = append(records, variant)
	}

	return &records, nil
}

// ReadProductVariant Find
func (m *Mongo) ReadProductVariant(dbname string, spu string, sku string) (*proto.VariantRecord, error) {
	vc := m.session.DB(dbname).C(variantsCollectionName)

	variant := Variant{}
	if err := vc.Find(bson.M{"spu": spu, "sku": sku}).One(&variant); err != nil {
		return nil, err
	}

	pricings := &proto.VariantRecord_Pricing{
		Current: variant.Pricings.Current,
	}
	images := []*proto.Image{}
	for _, ii := range variant.Images {
		image := &proto.Image{
			Type:  ii.Type,
			Paths: ii.Paths,
		}
		images = append(images, image)
	}
	record := proto.VariantRecord{
		Sku:              variant.Sku,
		Name:             variant.Name,
		Pricings:         pricings,
		Tracked:          variant.Tracked,
		ShippingCategory: variant.ShippingCategory,
		OptionValues:     variant.OptionValues,
		Stock:            variant.Stock,
		Images:           images,
		UpdatedAt:        variant.UpdatedAt,
		CreatedAt:        variant.CreatedAt,
	}

	return &record, nil
}

// UpdateProductVariant Find
func (m *Mongo) UpdateProductVariant(dbname string, spu string, sku string, record *proto.VariantRecord) error {
	vc := m.session.DB(dbname).C(variantsCollectionName)

	selector := bson.M{"spu": spu, "sku": sku}
	update := bson.M{"$set": bson.M{
		"name": record.Name,
		"pricings.current": record.Pricings.Current,
		"tracked": record.Tracked,
		"shippingCategory": record.ShippingCategory,
		"optionValues": record.OptionValues,
		"images": record.Images,
		"updated_at": time.Now().Unix(),
	}}

	return vc.Update(selector, update)
}

// DeleteProductVariant Remove
func (m *Mongo) DeleteProductVariant(dbname string, spu string, sku string) error {
	vc := m.session.DB(dbname).C(variantsCollectionName)
	selector := bson.M{"spu": spu, "sku": sku}
	return vc.Remove(selector)
}


// CreateProductOption Update
func (m *Mongo) CreateProductOption(dbname string, spu string, record *proto.ProductOption) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}
	option := Product_Option{
		Code:  record.Code,
		Value: record.Value,
	}
	updataData := bson.M{"$push": bson.M{
		"options": option}}

	return c.Update(selector, updataData)
}

// ReadProductOptions Find
func (m *Mongo) ReadProductOptions(dbname string, spu string) (*[]*proto.ProductRecord_OptionRecord, error) {

	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return nil, err
	}

	options := []*proto.ProductRecord_OptionRecord{}
	for _, io := range product.Options {
		if op, err := m.ReadOption(dbname, io.Code); err == nil {
			poos := []*proto.ProductRecord_OptionRecord_OptionValue{}
			for _, o := range op.Options {
				poo := &proto.ProductRecord_OptionRecord_OptionValue{
					Value:       o.Value,
					Description: o.Description,
				}
				poos = append(poos, poo)
			}
			option := &proto.ProductRecord_OptionRecord{
				Code:    io.Code,
				Value:   io.Code,
				Options: poos,
			}
			options = append(options, option)
		}
	}
	
	return &options, nil
}

// UpdateProductOption Update
func (m *Mongo) UpdateProductOption(dbname string, spu string, record *proto.ProductOption) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Options {
		if it.Code == record.Code {
			pullData := bson.M{"$pull": bson.M{
				"options": it}}
			if err := c.Update(selector, pullData); err != nil {
				return err
			}

			option := Product_Option{
				Code:  record.Code,
				Value: record.Value,
			}
			pushData := bson.M{"$push": bson.M{
				"options": option}}

			return c.Update(selector, pushData)
		}
	}

	return errors.New("没有找到 Code:" + record.Code)
}

// DeleteProductOption Update
func (m *Mongo) DeleteProductOption(dbname string, spu string, code string) error {
	c := m.session.DB(dbname).C(productsCollectionName)

	selector := bson.M{"spu": spu}

	product := Product{}
	if err := c.Find(selector).One(&product); err != nil {
		return err
	}

	for _, it := range product.Options {
		if it.Code == code {
			pullData := bson.M{"$pull": bson.M{
				"options": it}}
			return c.Update(selector, pullData)
		}
	}

	return errors.New("没有找到 Code:" + code)
}