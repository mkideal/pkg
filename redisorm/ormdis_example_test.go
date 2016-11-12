package orm

import (
	"github.com/mkideal/log"
)

//-----------
// sku table
//-----------

type SkuMeta struct {
	F_id    string
	F_name  string
	F_price string
}

func (SkuMeta) Name() string     { return "sku" }
func (SkuMeta) Fields() []string { return _sku_fields }

var skuMeta = &SkuMeta{
	F_id:    "id",
	F_name:  "name",
	F_price: "price",
}

var _sku_fields = []string{
	skuMeta.F_name,
	skuMeta.F_price,
}

type Sku struct {
	Id    int64
	Name  string
	Price int
}

func (sku Sku) Meta() TableMeta  { return skuMeta }
func (sku Sku) Key() interface{} { return sku.Id }

func (sku Sku) GetField(field string) interface{} {
	switch field {
	case skuMeta.F_id:
		return sku.Id
	case skuMeta.F_name:
		return sku.Name
	case skuMeta.F_price:
		return sku.Price
	}
	return nil
}

func (sku *Sku) SetField(field, value string) error {
	switch field {
	case skuMeta.F_name:
		sku.Name = value
	case skuMeta.F_price:
		return setInt(&sku.Price, value)
	}
	return nil
}

//---------------
// product table
//---------------

type ProductMeta struct {
	F_id    string
	F_name  string
	F_desc  string
	F_image string
	F_skuId string
}

func (ProductMeta) Name() string     { return "product" }
func (ProductMeta) Fields() []string { return _product_fields }

var productMeta = ProductMeta{
	F_name:  "name",
	F_desc:  "desc",
	F_image: "image",
	F_skuId: "sku_id",
}

var _product_fields = []string{
	productMeta.F_id,
	productMeta.F_name,
	productMeta.F_desc,
	productMeta.F_image,
	productMeta.F_skuId,
}

type Product struct {
	Id    int64
	Name  string
	Desc  string
	Image string
	SkuId int64
}

func (product Product) Key() interface{} { return product.Id }
func (product Product) Meta() TableMeta  { return productMeta }

func (product Product) GetField(field string) interface{} {
	switch field {
	case productMeta.F_name:
		return product.Name
	case productMeta.F_desc:
		return product.Desc
	case productMeta.F_image:
		return product.Image
	case productMeta.F_skuId:
		return product.SkuId
	}
	return nil
}

func (product *Product) SetField(field, value string) error {
	switch field {
	case productMeta.F_name:
		product.Name = value
	case productMeta.F_desc:
		product.Desc = value
	case productMeta.F_image:
		product.Image = value
	case productMeta.F_skuId:
		return setInt64(&product.SkuId, value)
	}
	return nil
}

// view: ProductDetail
type ViewProuctDetail struct{}

var _view_prouct_detail_fields = []string{
	productMeta.F_name,
	productMeta.F_desc,
	productMeta.F_image,
	productMeta.F_skuId,
}

func (ViewProuctDetail) Table() string     { return productMeta.Name() }
func (ViewProuctDetail) Fields() FieldList { return FieldSlice(_view_prouct_detail_fields) }

func Example_Basic() {
	defer log.Uninit(log.InitConsole(log.LvFATAL))

	engine := NewEngine("test", redisc)
	eng := engine.Core()
	eng.SetErrorHandler(func(action string, err error) error {
		log.Printf(ErrorHandlerDepth, log.LvWARN, "<%s>: %v", action, err)
		return err
	})
}

func Example_View() {
}
