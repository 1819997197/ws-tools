package model

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"regexp"
	"strings"
)

//MysqlTable mysql表格结构体
type MysqlTable struct {
	Table               string                 //表名 mysql表名
	TableName           string                 //表名 mysql 表中文名
	TableAlias          string                 //别名 用于结构体名称
	SQLCreate           string                 //建表语句
	Fields              map[string]SchemaField //字段集 [字段名]
	FieldNames          []string               //字段名集合 带``的字段名
	Indexes             map[string]SchemaIndex //索引
	HasPrimary          bool                   //是否有主键字段
	IsOnlyPrimary       bool                   //是否唯一主键
	OnlyPrimaryKeyField string                 //唯一主键字段
	Doc                 string                 //文档
	FieldStructDoc      bool                   //字段结构体文档
	IsModel             bool                   //是否具备建模条件
	HasTime             bool                   //是否有datetime date字段
	HasExtendType       bool                   //是否有指定Typexyz支持的扩展类型的ArrayString Integers Boolean,Timestamp类型
	HasState            bool                   //是否有状态字段
	StateFieldName      string                 //状态字段名
	HasCreated          bool                   //是否有创建时间字段
	CreatedFieldName    string                 // 创建字段的字段名
	HasUpdated          bool                   //是否有更新时间字段
	UpdatedFieldName    string                 // 更新字段的字段名
	HasDeleted          bool                   //是否有删除时间字段
	DeletedFieldName    string                 // 删除字段的字段名
	SearchDisabled      bool                   //不需要生成搜索代码
	PagingDisabled      bool                   //不需要生成分页代码
}

//FieldName2SQL 输出所有字段的拼接，用于 select 字段集 SQL语句
func (m MysqlTable) FieldName2SQL() string {
	return strings.Join(m.FieldNames, ", ")
}

// FieldAlias2Code 输出用于对应字段的赋值
func (m MysqlTable) FieldAlias2Code() string {
	var tmp string
	for _, fn := range m.FieldNames {
		f, _ := m.Field(fn)
		if tmp != "" {
			tmp += ", "
		}
		tmp += f.FieldAlias
	}
	return tmp
}

// Field 通过字段名获取字段信息，返回字段的结构体 SchemaField
func (m MysqlTable) Field(fieldName string) (schemaField SchemaField, err error) {
	if v, ok := m.Fields[fieldName]; ok {
		return v, nil
	}
	return SchemaField{}, errors.New("未找到匹配的字段信息")
}

// Primary2Code 获取主键的SQL条件语句和参数代码
func (m MysqlTable) Primary2Code() (whereSQL, paramInitCode, paramListCode string) {
	//b, _ := json.Marshal(m)
	//fmt.Println("json: ", string(b))
	if v, ok := m.Indexes["PRIMARY"]; ok {
		for _, fn := range v.FieldName {
			if whereSQL != "" {
				whereSQL += " AND "
			}
			if paramInitCode != "" {
				paramInitCode += ", "
				paramListCode += ", "
			}
			if v, ok := m.Fields[fn]; ok {
				paramInitCode += v.FieldAlias + " " + v.FieldType
				paramListCode += v.FieldAlias
			}
			whereSQL = whereSQL + fn + "=?"
		}
		if whereSQL != "" {
			whereSQL = "" + whereSQL
		}
	}
	return
}

//Tables 得到所有的表信息
func Tables(db *sqlx.DB, destTable []string, modelInit *ModelInit) (tables []*MysqlTable, err error) {
	var tableName string
	var arrTable []string
	arrTable, err = GetTableList(db) //获取所有的表
	if err != nil {
		return
	}

	// 获取所有的主键字段及类型
	for _, tableName = range arrTable {
		if len(destTable) > 0 {
			if !StringArrayExists(destTable, tableName) {
				continue
			}
		}
		var schema MysqlTable
		schema.Table = "`" + tableName + "`"
		// 获取字段
		var mysqlFullFields []MysqlTableField
		mysqlFullFields, err = GetTableColumn(db, tableName)
		if err != nil {
			return
		}
		var schemaFields = make(map[string]SchemaField)
		for _, f := range mysqlFullFields {
			var schemaField SchemaField
			schemaField = f.CommentX()
			if schemaField.IsTime {
				schema.HasTime = true
			}
			if schemaField.IsExtendType {
				schema.HasExtendType = true
			}
			if modelInit.IsState(schemaField.FieldAlias) {
				schema.HasState = true
				schema.StateFieldName = schemaField.SQLFieldName
			}
			if modelInit.IsCreated(schemaField.FieldAlias) {
				schema.HasCreated = true
				schema.CreatedFieldName = schemaField.SQLFieldName
			}
			if modelInit.IsUpdated(schemaField.FieldAlias) {
				schema.HasUpdated = true
				schema.UpdatedFieldName = schemaField.SQLFieldName
			}
			if modelInit.IsDeleted(schemaField.FieldAlias) {
				schema.HasDeleted = true
				schema.DeletedFieldName = schemaField.SQLFieldName
			}
			if schema.HasDeleted && schema.HasState {
				schema.IsModel = true
			}
			if schemaField.IsPrimary {
				schema.HasPrimary = true
			}
			schemaFields[f.FieldName()] = schemaField
			schema.FieldNames = append(schema.FieldNames, f.FieldName())
		}
		schema.Fields = schemaFields
		// 获取字段结束
		// 获取索引
		schema.Indexes = make(map[string]SchemaIndex)
		schema.Indexes, schema.IsOnlyPrimary, schema.OnlyPrimaryKeyField, err = GetTableIndexes(db, tableName)
		if err != nil {
			log.Print("GetTableIndexes err: ", err)
			return
		}
		// 获取索引结束
		err = GetTableExtInfo(db, &schema)
		if err != nil {
			return
		}
		tables = append(tables, &schema)
	}
	return
}

// GetTableList 获取数据库所有的表名
func GetTableList(db *sqlx.DB) (arrTable []string, err error) {
	var sqlQuery = "show tables"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		err = errors.New("show tables:" + err.Error())
		return
	}
	defer rows.Close()
	var tableName string
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			err = errors.New("show tables:" + err.Error())
			return
		}
		arrTable = append(arrTable, tableName)
	}
	return arrTable, nil
}

// GetTableColumn 从数据库获取指定的表的所有字段信息
func GetTableColumn(db *sqlx.DB, tableName string) (mysqlFullFields []MysqlTableField, err error) {
	err = db.Select(&mysqlFullFields, "show full columns from `"+tableName+"`")
	if err != nil {
		err = errors.New("show full columns: " + err.Error())
		return
	}
	return mysqlFullFields, nil
}

// GetTableExtInfo 获取表的中文别名，及对表别名，并获得创建的语句
// Alias 表别名用于定义结构体名称 TableName 中文名称用于说明 SQLCreate 建表SQL
func GetTableExtInfo(db *sqlx.DB, mysqlTable *MysqlTable) (err error) {
	var sqlQuery = "show create table " + mysqlTable.Table
	err = db.QueryRow(sqlQuery).Scan(&mysqlTable.Table, &mysqlTable.SQLCreate)
	if err != nil {
		return errors.New("show create table:" + err.Error())
	}

	mysqlTable.SQLCreate = strings.Replace(mysqlTable.SQLCreate, "，", ",", -1)
	var regSQLCreate = regexp.MustCompile(`AUTO_INCREMENT=\d+`)
	mysqlTable.SQLCreate = regSQLCreate.ReplaceAllString(mysqlTable.SQLCreate, "AUTO_INCREMENT=0")
	tmpArr := strings.SplitAfter(mysqlTable.SQLCreate, "ENGINE=")
	if len(tmpArr) != 2 {
		return nil
	}
	//表的备注规则 结构体名称,名称
	//例 User,用户表,SearchDisabled,PagingDisabled
	var reg = regexp.MustCompile(`COMMENT='(.*)'$`)
	match := reg.FindStringSubmatch(tmpArr[1])
	if len(match) == 2 {
		tableComment := match[1]
		tableComment = strings.Replace(tableComment, " ", "", -1)
		var reg = regexp.MustCompile(`([,]+)`)
		tableComment = reg.ReplaceAllString(tableComment, ",")
		//规则： 别名或表名称
		//规则2
		if strings.Index(tableComment, ",") > 0 {
			tmp := strings.Split(tableComment, ",")
			for _, v := range tmp {
				switch {
				case v == "SearchDisabled":
					mysqlTable.SearchDisabled = true
				case v == "PagingDisabled":
					mysqlTable.PagingDisabled = true
				case isAlias(v) == true:
					mysqlTable.TableAlias = v
				default:
					mysqlTable.TableName = v
				}
			}
			return nil
		}
	}
	if mysqlTable.TableAlias == "" {
		tmp := strings.Split(mysqlTable.Table, "_")
		for _, str := range tmp {
			mysqlTable.TableAlias += strings.Title(str)
		}
	}
	return nil
}

// IsAlias 判断是否是别名 只能是大驼峰命名法，例 User100 或UserId UserID
func isAlias(Alias string) bool {
	var reg = regexp.MustCompile(`^([A-Z]+)([A-Z0-9a-z_]+)$`)
	if reg.Match([]byte(Alias)) == false {
		return false
	}
	return true
}
