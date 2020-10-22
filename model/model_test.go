package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ProductData struct {
}

type Product struct {
	Id        int    `db:"primary_key"`
	FieldName string `db:"column:column_name_test"`
}

func (p *Product) TableName() string {
	return "product"
}

func TestParseStruct(t *testing.T) {
	product, _ := ParseStruct(Product{})
	productData, _ := ParseStruct(ProductData{})

	assert.Equal(t, "product", product.TableName)
	assert.Equal(t, "product_data", productData.TableName)

	assert.Equal(t, "Id", product.FieldAttrMap["id"])
	assert.Equal(t, "FieldName", product.FieldAttrMap["column_name_test"])
	assert.Equal(t, "id", product.PrimaryKey)
}
