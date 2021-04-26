package model

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

// SQLCreatorByIndexes 生成主键的条件语句 例 WHERE (`id` = ?) and (`pid` = ?)
func SQLCreatorByIndexes(schemaIndexes map[string]SchemaIndex, key string) string {
	var query string
	if v, ok := schemaIndexes[key]; ok {
		for _, fn := range v.FieldName {
			if query != "" {
				query += " AND "
			}
			query = query + "(" + fn + "=?)"
		}
		if query != "" {
			query = "WHERE " + query
		}
	}
	return query
}

// MySQLTableIndex 用于获取 show index from t 的结果
type MySQLTableIndex struct {
	Table        string         `db:"Table"`
	NonUnique    int8           `db:"Non_unique"` // 数据说明： 1 为 Normal 0 则为非Normal
	KeyName      string         `db:"Key_name"`
	SeqInIndex   int8           `db:"Seq_in_index"`
	ColumnName   string         `db:"Column_name"`
	Collation    sql.NullString `db:"Collation"`
	Cardinality  sql.NullString `db:"Cardinality"`
	SubPart      sql.NullString `db:"Sub_part"`
	Packed       sql.NullString `db:"Packed"`
	Null         sql.NullString `db:"Null"`
	IndexType    sql.NullString `db:"Index_type"`
	Comment      sql.NullString `db:"Comment"`
	IndexComment sql.NullString `db:"Index_comment"`
	Visible      string         `db:"Visible"`
	Expression   sql.NullTime   `db:"Expression"`
}

// SchemaIndex 表的索引结构体
type SchemaIndex struct {
	Name      string
	FieldName []string
	IsPrimary bool // 是否主键索引
	IsUnique  bool // 是否Unique索引
}

// GetTableIndexes 从数据库获取指定的表的所有索引信息
func GetTableIndexes(db *sqlx.DB, tableName string) (schemaIndexes map[string]SchemaIndex, IsOnlyPrimary bool, OnlyPrimaryFieldName string, err error) {
	var MySQLTableIndexes []MySQLTableIndex
	err = db.Select(&MySQLTableIndexes, "show index from `"+tableName+"`")
	if err != nil {
		err = errors.New("GetTableIndexes : " + err.Error())
		return
	}
	var schemaIndex SchemaIndex
	schemaIndexes = make(map[string]SchemaIndex)
	indexLen := len(MySQLTableIndexes) - 1
	for index, f := range MySQLTableIndexes {
		if schemaIndex.Name != f.KeyName {
			if schemaIndex.Name != "" {
				schemaIndexes[schemaIndex.Name] = schemaIndex
			}
			schemaIndex = SchemaIndex{}
		}
		schemaIndex.Name = f.KeyName
		schemaIndex.FieldName = append(schemaIndex.FieldName, "`"+f.ColumnName+"`")
		if f.NonUnique == 0 {
			if f.KeyName == "PRIMARY" {
				schemaIndex.IsPrimary = true
			} else {
				schemaIndex.IsUnique = true
			}
		}
		if index == indexLen {
			schemaIndexes[schemaIndex.Name] = schemaIndex
		}
	}
	if v, ok := schemaIndexes["PRIMARY"]; ok {
		if len(v.FieldName) == 1 {
			IsOnlyPrimary = true
			OnlyPrimaryFieldName = v.FieldName[0]
		}
	}
	return schemaIndexes, IsOnlyPrimary, OnlyPrimaryFieldName, nil
}
