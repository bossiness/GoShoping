package db

import (
	proto "btdxcx.com/micro/product-srv/proto/product"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Attribute
	Option
	Product
}

var (
	db DB
)

// Register db Imp
func Register(backend DB) {
	db = backend
}

// Init 数据库初始化
func Init() error {
	return db.Init()
}

// Deinit 析构
func Deinit() {
	db.Deinit()
}

// Attribute is Attribute数据接口
type Attribute interface {
	CreateAttribute(string, *proto.AttributesRecord) error
	ReadAttributes(string, int, int) (*[]*proto.AttributesRecord, error)
	ReadAttribute(string, string) (*proto.AttributesRecord, error)
	UpdateAttribute(string, string, *proto.AttributesRecord) error
	DeleteAttribute(string, string) error
}

// CreateAttribute create
func CreateAttribute(dbname string, record *proto.AttributesRecord) error {
	return db.CreateAttribute(dbname, record)
}

// ReadAttributes read list
func ReadAttributes(dbname string, offset int, limit int) (*[]*proto.AttributesRecord, error) {
	return db.ReadAttributes(dbname, offset, limit)
}

// ReadAttribute read
func ReadAttribute(dbname string, code string) (*proto.AttributesRecord, error) {
	return db.ReadAttribute(dbname, code)
}

// UpdateAttribute update
func UpdateAttribute(dbname string, code string, record *proto.AttributesRecord) error {
	return db.UpdateAttribute(dbname, code, record)
}

// DeleteAttribute delete
func DeleteAttribute(dbname string, code string) error {
	return db.DeleteAttribute(dbname, code)
}

// Option is Option数据接口
type Option interface {
	CreateOption(string, *proto.OptionRecord) error
	ReadOptions(string, int, int) (*[]*proto.OptionRecord, error)
	ReadOption(string, string) (*proto.OptionRecord, error)
	UpdateOption(string, string, *proto.OptionRecord) error
	DeleteOption(string, string) error
}

// CreateOption create
func CreateOption(dbname string, record *proto.OptionRecord) error {
	return db.CreateOption(dbname, record)
}

// ReadOptions read
func ReadOptions(dbname string, offset int, limit int) (*[]*proto.OptionRecord, error) {
	return db.ReadOptions(dbname, offset, limit)
}

// ReadOption read
func ReadOption(dbname string, code string) (*proto.OptionRecord, error) {
	return db.ReadOption(dbname, code)
}

// UpdateOption update
func UpdateOption(dbname string, code string, record *proto.OptionRecord) error {
	return db.UpdateOption(dbname, code, record)
}

// DeleteOption delete
func DeleteOption(dbname string, code string) error {
	return db.DeleteOption(dbname, code)
}

// Product is Product数据接口
type Product interface {
	CreateProduct(dbname string, record *proto.ProductRecord) error
	ReadProducts(dbname string, offset int, limit int) (*[]*proto.ProductRecord, error)
	ReadProduct(dbname string, spu string) (*proto.ProductRecord, error)
	UpdateProduct(dbname string, spu string, name string, description string) error
	DeleteProduct(dbname string, spu string) error
	TaxonProducts(dbname string, taxonCode string, offset int, limit int) (*[]*proto.ProductRecord, error)
	UpdateProductTaxons(dbname string, spu string, main string, others []string) error 
	CreateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error
	ReadProductAttributes(dbname string, spu string) (*[]*proto.ProductRecord_AttributesRecord, error)
	UpdateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error
	DeleteProductAttribute(dbname string, spu string, code string) error
	CreateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error
	UpdateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error
	DeleteProductAssociation(dbname string, spu string, code string) error
	CreateProductImage(dbname string, spu string, record *proto.Image) error
	UpdateProductImage(dbname string, spu string, record *proto.Image) error
	DeleteProductImage(dbname string, spu string, itype string) error
	CreateProductReview(dbname string, spu string, record *proto.ReviewsRecord) error
	ReadProductReviews(dbname string, spu string, offset int, limit int) (*[]*proto.ReviewsRecord, error)
	ReadProductReview(dbname string, id string) (*proto.ReviewsRecord, error)
	UpdateProductReview(dbname string, record *proto.ReviewsRecord) error
	DeleteProductReview(dbname string, id string) error
	AcceptProductReview(dbname string, id string) error
	RejectProductReview(dbname string, spu string, id string) error
	CreateProductVariant(dbname string, spu string, record *proto.VariantRecord) error
	ReadProductVariants(dbname string, spu string, offset int, limit int) (*[]*proto.VariantRecord, error)
	ReadProductVariant(dbname string, spu string, sku string) (*proto.VariantRecord, error)
	UpdateProductVariant(dbname string, spu string, sku string, record *proto.VariantRecord) error 
	DeleteProductVariant(dbname string, spu string, sku string) error
	CreateProductOption(dbname string, spu string, record *proto.ProductOption) error
	ReadProductOptions(dbname string, spu string) (*[]*proto.ProductRecord_OptionRecord, error)
	UpdateProductOption(dbname string, spu string, record *proto.ProductOption) error
	DeleteProductOption(dbname string, spu string, code string) error
	EnableProduct(dbname string, spu string, enabled bool) error
}

// CreateProduct Insert
func CreateProduct(dbname string, record *proto.ProductRecord) error {
	return db.CreateProduct(dbname, record)
}

// ReadProducts Find
func ReadProducts(dbname string, offset int, limit int) (*[]*proto.ProductRecord, error) {
	return db.ReadProducts(dbname, offset, limit)
}

// ReadProduct Find
func ReadProduct(dbname string, spu string) (*proto.ProductRecord, error) {
	return db.ReadProduct(dbname, spu)
}


// UpdateProduct Update
func UpdateProduct(dbname string, spu string, name string, description string) error {
	return db.UpdateProduct(dbname, spu, name, description)
}

// DeleteProduct Remove
func DeleteProduct(dbname string, spu string) error {
	return db.DeleteProduct(dbname, spu)
}

// TaxonProducts taxon products
func TaxonProducts(dbname string, taxonCode string, offset int, limit int) (*[]*proto.ProductRecord, error) {
	return db.TaxonProducts(dbname, taxonCode, offset, limit)
}

// UpdateProductTaxons Update
func UpdateProductTaxons(dbname string, spu string, main string, others []string) error {
	return db.UpdateProductTaxons(dbname, spu, main, others)
}

// CreateProductAttribute Update
func CreateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error {
	return db.CreateProductAttribute(dbname, spu, record)
}

func ReadProductAttributes(dbname string, spu string) (*[]*proto.ProductRecord_AttributesRecord, error) {
	return db.ReadProductAttributes(dbname, spu)
}

// UpdateProductAttribute Update
func UpdateProductAttribute(dbname string, spu string, record *proto.ProductAttribute) error {
	return db.UpdateProductAttribute(dbname, spu, record)
}

// DeleteProductAttribute Update
func DeleteProductAttribute(dbname string, spu string, code string) error {
	return db.DeleteProductAttribute(dbname, spu, code)
}

// CreateProductAssociation Update
func CreateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error {
	return db.CreateProductAssociation(dbname, spu, record)
}

// UpdateProductAssociation Update
func UpdateProductAssociation(dbname string, spu string, record *proto.ProductAssociation) error {
	return db.UpdateProductAssociation(dbname, spu, record)
}

// DeleteProductAssociation Update
func DeleteProductAssociation(dbname string, spu string, code string) error {
	return db.DeleteProductAssociation(dbname, spu, code)
}

// CreateProductImage Update
func CreateProductImage(dbname string, spu string, record *proto.Image) error {
	return db.CreateProductImage(dbname, spu, record)
}

// UpdateProductImage Update
func UpdateProductImage(dbname string, spu string, record *proto.Image) error {
	return db.UpdateProductImage(dbname, spu, record)
}

// DeleteProductImage Update
func DeleteProductImage(dbname string, spu string, itype string) error {
	return db.DeleteProductImage(dbname, spu, itype)
}

// CreateProductReview Update
func CreateProductReview(dbname string, spu string, record *proto.ReviewsRecord) error {
	return db.CreateProductReview(dbname, spu, record)
}

// ReadProductReviews Find
func ReadProductReviews(dbname string, spu string, offset int, limit int) (*[]*proto.ReviewsRecord, error) {
	return db.ReadProductReviews(dbname, spu, offset, limit)
}

// ReadProductReview Find
func ReadProductReview(dbname string, id string) (*proto.ReviewsRecord, error) {
	return db.ReadProductReview(dbname, id)
}

// UpdateProductReview Update
func UpdateProductReview(dbname string, record *proto.ReviewsRecord) error {
	return db.UpdateProductReview(dbname, record)
}

// DeleteProductReview Remove id
func DeleteProductReview(dbname string, id string) error {
	return db.DeleteProductReview(dbname, id)
}

// AcceptProductReview Update
func AcceptProductReview(dbname string, id string) error {
	return db.AcceptProductReview(dbname, id)
}

// RejectProductReview Update
func RejectProductReview(dbname string, spu string, id string) error {
	return db.RejectProductReview(dbname, spu, id)
}

// CreateProductVariant Insert
func CreateProductVariant(dbname string, spu string, record *proto.VariantRecord) error {
	return db.CreateProductVariant(dbname, spu, record)
}

// ReadProductVariants Find
func ReadProductVariants(dbname string, spu string, offset int, limit int) (*[]*proto.VariantRecord, error) {
	return db.ReadProductVariants(dbname, spu, offset, limit)
}

// ReadProductVariant Find
func ReadProductVariant(dbname string, spu string, sku string) (*proto.VariantRecord, error) {
	return db.ReadProductVariant(dbname, spu, sku)
}

// UpdateProductVariant Find
func UpdateProductVariant(dbname string, spu string, sku string, record *proto.VariantRecord) error {
	return db.UpdateProductVariant(dbname, spu, sku, record)
}

// DeleteProductVariant Remove
func DeleteProductVariant(dbname string, spu string, sku string) error {
	return db.DeleteProductVariant(dbname, spu, sku)
}

// CreateProductOption Update
func CreateProductOption(dbname string, spu string, record *proto.ProductOption) error {
	return db.CreateProductOption(dbname, spu, record)
}

// ReadProductOptions Find
func ReadProductOptions(dbname string, spu string) (*[]*proto.ProductRecord_OptionRecord, error) {
	return db.ReadProductOptions(dbname, spu)
}

// UpdateProductOption Update
func UpdateProductOption(dbname string, spu string, record *proto.ProductOption) error {
	return db.UpdateProductOption(dbname, spu, record)
}

// DeleteProductOption Update
func DeleteProductOption(dbname string, spu string, code string) error {
	return db.DeleteProductOption(dbname, spu, code)
}

// EnableProduct Update
func EnableProduct(dbname string, spu string, enabled bool) error {
	return db.EnableProduct(dbname, spu, enabled)
}