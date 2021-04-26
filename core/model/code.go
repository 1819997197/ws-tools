package model

import (
	"strings"
)

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
	_, paramInitCode, _ := table.Primary2Code()
	strMethod += "//go:generate mockgen -destination=./mock/" + table.Table + ".go -package=mock {{" + packageName + "}} " + table.TableAlias + "RepositoryIFace\n"
	strMethod += "type " + table.TableAlias + "RepositoryIFace interface {\n"
	strMethod += "\tFindInfoList(fields string, vars ...string) ([]*" + table.TableAlias + "Model, error)\n"
	strMethod += "\tFindInfoByPrimaryID(fields string, " + paramInitCode + ") (*" + table.TableAlias + "Model, error)\n"
	strMethod += "\tCreate(data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tTxCreate(tx *gorm.DB, data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tSave(data *" + table.TableAlias + "Model) error\n"
	strMethod += "\tTxSave(tx *gorm.DB, data *" + table.TableAlias + "Model) error\n"
	strMethod += "}\n\n"
	return strMethod
}

// GetStructMethod 获取表结构
func GetStructMethod(table *MysqlTable) string {
	var strMethod string
	whereSQL, paramInitCode, paramListCode := table.Primary2Code()
	strMethod += "\n// FindInfoList 根据条件查找列表\n"
	strMethod += "func (entity *" + table.TableAlias + "Entity) FindInfoList(fields string, vars ...string) ([]*" + table.TableAlias + "Model, error) { \n"
	strMethod += "\tvar sqlWhere = \"1=1\"\n"
	strMethod += "\tif len(vars) > 0 {\n"
	strMethod += "\t\tsqlWhere = vars[0]\n"
	strMethod += "\t}\n"
	strMethod += "\tvar list []*" + table.TableAlias + "Model\n"
	strMethod += "\tquery := entity.db.Where(sqlWhere)\n"
	strMethod += "\tif fields != \"\" {\n"
	strMethod += "\t\tquery = query.Select(fields)\n"
	strMethod += "\t} else {\n"
	strMethod += "\t\tquery = query.Select(" + strings.Join(table.FieldNames, ", ") + ")\n"
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

	if table.HasPrimary {
		strMethod += "\n// FindInfoByPrimaryID 根据主键查找一条记录\n"
		strMethod += "func (entity *" + table.TableAlias + "Entity) FindInfoByPrimaryID(fields string, " + paramInitCode + ") (*" + table.TableAlias + "Model, error) { \n"
		strMethod += "\tvar model = &" + table.TableAlias + "Model{}\n"
		strMethod += "\tquery := entity.db.Where(\"" + whereSQL + "\", " + paramListCode + ")\n"
		strMethod += "\tif fields != \"\" {\n"
		strMethod += "\t\tquery = query.Select(fields)\n"
		strMethod += "\t} else {\n"
		strMethod += "\t\tquery = query.Select(" + strings.Join(table.FieldNames, ", ") + ")\n"
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
	}

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
	return strMethod
}