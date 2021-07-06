package model

// Header 生成代码
func Header(table *MysqlTable, packageName string) string {
	var fileContent = ""
	fileContent += "package " + packageName + "\n\n"
	fileContent += "import (\n"
	fileContent += "\t\"git.wondershare.cn/DCStudio/chaos_go/core/database\"\n"
	fileContent += "\t\"github.com/jinzhu/gorm\"\n"
	if table.HasTime {
		fileContent += "\t\"time\"\n"
	}
	fileContent += ")\n"
	fileContent += "\n"
	return fileContent
}

/*
 * 生成表的代码
 * structDescription: 结构体描述
 * structDoc: 结构体定义
 * tableNameDoc: 表名
 * entityNameDoc: 实体定义
 * entityDoc: 实体对象
 */
func TableDoc(table *MysqlTable, packageName string) string {
	var structDoc = ""
	for _, fieldname := range table.FieldNames {
		if field, ok := table.Fields[fieldname]; ok {
			structDoc += "\t" + field.FieldAlias + " "
			structDoc += "\t" + field.FieldType
			structDoc += "\t`json:\"" + field.FieldName + "\"`"
			structDoc += "\t// " + field.FieldTitle + " 类型: " + field.FieldType
			if field.IsPrimary {
				structDoc += " 主健字段（Primary Key）"
			}
			if field.IsAutoIncrement {
				structDoc += " 自增长字段 "
			}
			if field.FieldDescription != "" {
				structDoc += " 说明: " + field.FieldDescription
			}
			if field.FieldDefault != "" {
				structDoc += " 默认值: " + field.FieldDefault
			}

			structDoc += "\n"
		}
	}
	structDoc = "type " + table.TableAlias + "Model struct {\n" + structDoc + "}\n"
	var structDescription string = "// " + table.TableAlias + "Model 针对数据库表 " + table.Table + " 的结构体定义\n"
	var tableNameDoc = "\nfunc (" + table.TableAlias + "Model) TableName() string {\n\treturn \"" + table.Table + "\"\n}\n"
	var entityNameDoc = "\ntype " + table.TableAlias + "Entity struct {\n\tdb *database.DB\n}\n"
	var entityDoc = "\nfunc New" + table.TableAlias + "Entity(db *database.DB) " + table.TableAlias + "RepositoryIFace {\n\treturn &" + table.TableAlias + "Entity{db}\n}\n"
	return GetInterfaceMethod(table, packageName) + structDescription + structDoc + tableNameDoc + entityNameDoc + entityDoc + GetStructMethod(table)
}

func GetInterfaceMethod(table *MysqlTable, packageName string) string {
	var strMethod string
	strMethod += "//go:generate mockgen -destination=./mock/" + table.Table + ".go -package=mock {{" + packageName + "}} " + table.TableAlias + "RepositoryIFace\n"
	strMethod += "type " + table.TableAlias + "RepositoryIFace interface {\n"
	strMethod += "\tCreate(data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tTxCreate(tx *gorm.DB, data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tSave(data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tTxSave(tx *gorm.DB, data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tUpdateFields(where, updateFiles map[string]interface{}) error\n"
	strMethod += "\tTxUpdateFields(tx *gorm.DB, where, updateFiles map[string]interface{}) error\n"
	strMethod += "}\n\n"
	return strMethod
}

// GetStructMethod 获取表结构
func GetStructMethod(table *MysqlTable) string {
	var strMethod string
	strMethod += "\n// getListBy 根据条件查找列表\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) getListBy(fields string, where map[string]interface{}) ([]*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tquery := entity.db.Where(where)\n"
	strMethod += "\treturn entity.findAll(query, fields)\n"
	strMethod += "}\n"

	strMethod += "\n// where: ws_id = ? AND m_pid IN(?) AND attribute_key=?  args:[]interface{}{12039, []int{846}, \"feature\"}\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) getListByWhere(fields string, where string, args []interface{}) ([]*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tquery := entity.db.Where(where, args...)\n"
	strMethod += "\treturn entity.findAll(query, fields)\n"
	strMethod += "}\n"

	strMethod += "\n// getBy 根据条件查找一条记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) getBy(fields string, where map[string]interface{}) (*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tquery := entity.db.Where(where)\n"
	strMethod += "\treturn entity.findOne(query, fields)\n"
	strMethod += "}\n"

	strMethod += "\n// where: ws_id = ? AND m_pid IN(?) AND attribute_key=?  args:[]interface{}{12039, []int{846}, \"feature\"}\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) getByWhere(fields string, where string, args []interface{}) (*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tquery := entity.db.Where(where, args...)\n"
	strMethod += "\treturn entity.findOne(query, fields)\n"
	strMethod += "}\n"

	strMethod += "func (entity *" + table.TableAlias + "Entity) findAll(query *gorm.DB, fields string) ([]*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tvar list []*" + table.TableAlias + "Model\n"
	strMethod += "\tif fields != \"\" {\n"
	strMethod += "\t\tquery = query.Select(fields)\n"
	strMethod += "\t}\n"
	strMethod += "\terr := query.Find(&list).Error\n"
	strMethod += "\tif err == gorm.ErrRecordNotFound {\n"
	strMethod += "\t\treturn nil, nil\n"
	strMethod += "\t}\n"
	strMethod += "\tif err != nil {\n"
	strMethod += "\t\treturn nil, err\n"
	strMethod += "\t}\n"
	strMethod += "\treturn list, nil\n"
	strMethod += "}\n"

	strMethod += "func (entity *" + table.TableAlias + "Entity) findOne(query *gorm.DB, fields string) (*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tvar model = &" + table.TableAlias + "Model{}\n"
	strMethod += "\tif fields != \"\" {\n"
	strMethod += "\t\tquery = query.Select(fields)\n"
	strMethod += "\t}\n"
	strMethod += "\terr := query.First(model).Error\n"
	strMethod += "\tif err == gorm.ErrRecordNotFound {\n"
	strMethod += "\t\treturn nil, nil\n"
	strMethod += "\t}\n"
	strMethod += "\tif err != nil {\n"
	strMethod += "\t\treturn nil, err\n"
	strMethod += "\t}\n"
	strMethod += "\treturn model, nil\n"
	strMethod += "}\n"

	strMethod += GetSaveFunc(table)
	strMethod += GetUpdateFunc(table)
	return strMethod
}

//GetSaveFunc 保存函数代码
func GetSaveFunc(table *MysqlTable) string {
	var strMethod string
	strMethod += "\n// Create 写入记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) Create(data *" + table.TableAlias + "Model) error {\n"
	strMethod += "\treturn entity.db.Create(data).Error\n"
	strMethod += "}\n"

	strMethod += "\n// TxCreate 写入记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) TxCreate(tx *gorm.DB, data *" + table.TableAlias + "Model) error {\n"
	strMethod += "\treturn tx.Create(data).Error\n"
	strMethod += "}\n"
	return strMethod
}

//GetUpdateFunc 更新函数
func GetUpdateFunc(table *MysqlTable) string {
	var strMethod string
	strMethod += "\n// Save 更新记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) Save(data *" + table.TableAlias + "Model) error {\n"
	strMethod += "\treturn entity.db.Save(data).Error\n"
	strMethod += "}\n"

	strMethod += "\n// TxSave 更新记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) TxSave(tx *gorm.DB, data *" + table.TableAlias + "Model) error {\n"
	strMethod += "\treturn tx.Save(data).Error\n"
	strMethod += "}\n"

	strMethod += "\n// UpdateFields 根据条件更新记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) UpdateFields(where, updateFiles map[string]interface{}) error {\n"
	strMethod += "\treturn entity.db.Model(" + table.TableAlias + "Model{}).Where(where).UpdateColumns(updateFiles).Error\n"
	strMethod += "}\n"

	strMethod += "\n// TxUpdateFields 根据条件更新记录\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) TxUpdateFields(tx *gorm.DB, where, updateFiles map[string]interface{}) error {\n"
	strMethod += "\treturn tx.Model(" + table.TableAlias + "Model{}).Where(where).UpdateColumns(updateFiles).Error\n"
	strMethod += "}\n"
	return strMethod
}
