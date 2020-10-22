package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wform/worm/model"
)

func TestSql_BuildInsert(t *testing.T) {
	objSql := getSqlObj()
	objSql.Rows([]map[string]interface{}{
		{
			"update_time": 15014245780,
		},
		{
			"update_time": 15014245781,
		},
	})
	sql := objSql.BuildInsert()
	assert.Equal(t, escapeSpace("INSERT INTO product(update_time) VALUES (15014245780),(15014245781)"), escapeSpace(sql))
}

func TestSql_BuildDelete(t *testing.T) {
	objSql := getSqlObj()
	objSql.Cond.Where("id=100")
	sql := objSql.BuildDelete()
	assert.Equal(t, escapeSpace("DELETE FROM product  WHERE id=100"), escapeSpace(sql))
}

func TestSql_BuildUpdate(t *testing.T) {
	objSql := getSqlObj()
	objSql.Cond.Where("id=100")
	objSql.Rows([]map[string]interface{}{
		{
			"update_time": 15014245785,
		},
	})
	sql := objSql.BuildUpdate()
	assert.Equal(t, escapeSpace("UPDATE product SET update_time=15014245785 WHERE id=100"), escapeSpace(sql))
}

func TestSql_BuildSelect(t *testing.T) {
	objSql := getSqlObj()
	objSql.Cond.Where("test id IN ?", []string{"test1", "test23"})
	sql := objSql.Select("id").Having(" total>1").Group(" id").Order("id desc").BuildSelect()
	assert.Equal(t, escapeSpace("SELECT id FROM product WHERE test id IN ('test1','test23') GROUP BY  id HAVING  total>1 ORDER BY id desc"), escapeSpace(sql))

	objSql = getSqlObj()
	objSql.Cond.Where("id=?", 100)
	sql = objSql.BuildSelect()
	assert.Equal(t, escapeSpace("SELECT * FROM product WHERE id='100'"), escapeSpace(sql))
}

func TestCondition_Where(t *testing.T) {
	objSql := getSqlObj()
	objSql.Cond.Where(Product{Id: 100})
	sql := objSql.BuildSelect()
	assert.Equal(t, escapeSpace("SELECT * FROM product WHERE id = 100"), escapeSpace(sql))

	objSql = getSqlObj()
	objSql.Cond.Not("id=?", 100)
	sql = objSql.BuildSelect()
	assert.Equal(t, escapeSpace("SELECT * FROM product WHERE  NOT ( id='100')"), escapeSpace(sql))
}

func escapeSpace(from string) (to string) {
	return strings.Replace(from, " ", "", -1)
}

// common object for test
func getSqlObj() *Sql {
	objSql := &Sql{}
	product := Product{}
	objSql = &Sql{}
	objSql.modelStruct, _ = model.ParseStruct(product)
	return objSql
}

type Product struct {
	Id        int    `worm:"primary_key"`
	FieldName string `worm:"column:field_name_test"`
}
