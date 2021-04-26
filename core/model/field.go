package model

import (
	"database/sql"
	"regexp"
	"strings"
)

//MysqlTableField 用来取读show full columns from 得到的mysql表字段结构体
type MysqlTableField struct {
	Field      string         `db:"Field"`
	Type       string         `db:"Type"`
	Collation  sql.NullString `db:"Collation"`
	Null       string         `db:"Null"`
	Key        string         `db:"Key"`
	Default    sql.NullString `db:"Default"`
	Extra      string         `db:"Extra"`
	Privileges string         `db:"Privileges"`
	Comment    string         `db:"Comment"`
}

// FieldAliasFromName 首字母转换成大写
func (f MysqlTableField) FieldAliasFromName() string {
	if strings.Index(f.Field, "_") > -1 {
		var arrStr = strings.Split(f.Field, "_")
		var alias string
		for _, v := range arrStr {
			alias += strings.Title(v)
		}
		return alias
	}
	return strings.Title(f.Field)
}

// FieldName 输出适合sql引用的字段名
func (f MysqlTableField) FieldName() string {
	return "`" + f.Field + "`"
}

// CommentX 通过字段备注规则获取设置,如果备注为空，则会从字段名，字段类型进行判断转换
func (f MysqlTableField) CommentX() (schemaField SchemaField) {
	//规则一：名字
	//规则二: 类型,名字
	//规则三: 类型,别名,名字
	//规则四: 类型,别名,名字，备注
	// 字段名，加``
	schemaField.FieldName = f.Field
	schemaField.SQLFieldName = "`" + f.Field + "`"
	//默认值
	schemaField.FieldDefault = f.Default.String

	schemaField.FieldDefault = strings.Replace(schemaField.FieldDefault, "'", "", -1)
	schemaField.FieldDefault = strings.Replace(schemaField.FieldDefault, " ", "", -1)
	if schemaField.FieldDefault == "NULL" {
		schemaField.FieldDefault = ""
	}

	schemaField.SQLType = f.Type
	f.Comment = strings.Replace(f.Comment, " ", "", -1)
	if f.Comment == "" {
		// 没有备注规则，进行转换 -- begin

		schemaField.FieldTitle = f.FieldAliasFromName()
		schemaField.FieldAlias = schemaField.FieldTitle
		var unsigned bool
		schemaField.FieldType, unsigned = f.FieldTypeX()
		if unsigned {
			schemaField.FieldType = "u" + schemaField.FieldType
		}

	} else {
		f.Comment = strings.Replace(f.Comment, "，", ",", -1)
		var regSQLCreate = regexp.MustCompile(`([,]+)`)
		f.Comment = regSQLCreate.ReplaceAllString(f.Comment, ",")
		var arrTmp = strings.Split(f.Comment, ",")
		if len(arrTmp) == 1 {
			//规则一,备注只有一项时，符合别名规则则为别名，符合类型则为类型，否则为字段标题名称
			ok, t := f.isFieldType(arrTmp[0])
			if ok {
				schemaField.FieldType = t
			} else {
				if f.isFieldAlias(arrTmp[0]) {
					schemaField.FieldAlias = arrTmp[0]
				} else {
					schemaField.FieldTitle = arrTmp[0]
				}
			}
		}
		if len(arrTmp) > 1 {
			for _, tmp := range arrTmp {
				ok, t := f.isFieldType(tmp)
				if ok {
					schemaField.FieldType = t
				} else {
					if f.isFieldAlias(tmp) && schemaField.FieldAlias == "" {
						schemaField.FieldAlias = tmp
					} else {
						if schemaField.FieldTitle == "" {
							schemaField.FieldTitle = tmp
						} else {
							schemaField.FieldDescription = tmp
						}
					}
				}
			}
		}
		if schemaField.FieldAlias == "" {
			schemaField.FieldAlias = f.FieldAliasFromName()
		}
		if schemaField.FieldType == "" {
			var unsigned bool
			schemaField.FieldType, unsigned = f.FieldTypeX()
			if unsigned && schemaField.FieldType != "float32" && schemaField.FieldType != "float64" {
				schemaField.FieldType = "u" + schemaField.FieldType
			}
		}
	}
	if strings.HasPrefix(schemaField.FieldType, "int") || strings.HasPrefix(schemaField.FieldType, "uint") {
		schemaField.IsInteger = true
	}
	if strings.HasPrefix(schemaField.FieldType, "time.Time") {
		schemaField.IsTime = true
	}
	if strings.HasPrefix(schemaField.FieldType, "typexyz.") {
		schemaField.IsExtendType = true
	}
	if f.Key == "PRI" {
		schemaField.IsPrimary = true
	}
	if f.Extra == "auto_increment" {
		schemaField.IsAutoIncrement = true
	}
	return
}

// FieldTypeX 拆解字段类型
func (f MysqlTableField) FieldTypeX() (t string, unsigned bool) {
	var regSQLCreate = regexp.MustCompile(`\(\d+\)`)
	f.Type = regSQLCreate.ReplaceAllString(f.Type, "")
	var tmpArr []string
	tmpArr = strings.Split(f.Type, " ")
	if len(tmpArr) == 0 {
		return
	}
	t = tmpArr[0]
	if len(tmpArr) == 2 {
		if tmpArr[1] == "unsigned" {
			unsigned = true
		}
	}
	if ok, tmp := f.isFieldType(t); ok {
		t = tmp
	}
	return
}

// IsFieldAlias 判断是否是别名 只能是大驼峰命名法，例 User100 或UserId UserID
func (f MysqlTableField) isFieldAlias(Alias string) bool {
	var reg = regexp.MustCompile(`^([A-Z]+)([A-Z0-9a-z_]+)$`)
	if reg.Match([]byte(Alias)) == false {
		return false
	}
	return true
}

//isUnsigned 是否无符号字段
func (f MysqlTableField) isUnsigned() bool {
	return true
}

// IsFieldType 判断是否支持的类型
func (f MysqlTableField) isFieldType(t string) (ok bool, goLangType string) {
	t = strings.ToLower(t)
	if _, ok = ExtendFieldTypeLimit[t]; ok {
		return ok, ExtendFieldTypeLimit[t]
	}
	return ok, ""
}

// SchemaField mysql字段结构体
type SchemaField struct {
	FieldName        string //不带``的字段名
	SQLFieldName     string //带``的字段名
	FieldAlias       string //结构体成员名称,如果没有设置别名将字段首字母变成大写。
	FieldTitle       string //中文名称
	SQLType          string // 字段类型--数据库类型
	FieldType        string //golang类型
	FieldDefault     string
	FieldDescription string //说明
	IsInteger        bool   //是否整型
	IsTime           bool   //是否日期时间类型
	IsExtendType     bool   //是否扩展类型
	IsAutoIncrement  bool   //是否自增长字段
	IsPrimary        bool   //是否主键字段
}

// Default 获取字段的默认值
func (sf SchemaField) Default() string {

	if sf.FieldDefault == "" {
		switch {
		case sf.IsInteger:
			return "0"
		case sf.FieldType == "time.Time":
			return "time.Now()"
		case sf.FieldType == "typexyz.ArrayString" || sf.FieldType == "typexyz.Integers":
			return sf.FieldType + "{}"
		case sf.FieldType == "typexyz.Boolean":
			return "false"
		case sf.FieldType == "typexyz.Timestamp":
			switch sf.FieldAlias {
			case "Deleted":
			case "Deprecated":
			case "Updated":
				return "typexyz.NewTimestamp(0)"
			case "Created":
				return "typexyz.Now()"
			default:
				return "typexyz.Now()"
			}
		}
	} else {
		switch {
		case sf.FieldType == "time.Time":
			return "time.Now()"
		case sf.FieldType == "typexyz.ArrayString" || sf.FieldType == "typexyz.Integers":
			return sf.FieldType + "{}"
		case sf.FieldType == "typexyz.Boolean":
			return "false"
		case sf.FieldType == "typexyz.Timestamp":
			switch sf.FieldAlias {
			case "Deleted", "Deprecated", "Updated":
				return "typexyz.NewTimestamp(0)"
			case "Created":
				return "typexyz.Now()"
			default:
				return "typexyz.Now()"
			}
		case sf.IsInteger:
			return sf.FieldDefault
		default:
			return "\"" + sf.FieldDefault + "\""
		}
	}
	return sf.FieldDefault
}
