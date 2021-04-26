package model

// TableMapInit 用于注册map数据
type TableMapInit struct {
	Name  string
	Key   string
	Value string
	Data  map[interface{}]interface{}
}

// ExtendFieldTypeLimit 扩展类型
var ExtendFieldTypeLimit = map[string]string{
	// ----- mysql ----
	"bigint":    "int64",
	"tinyint":   "int8",
	"int":       "int32",
	"smallint":  "int32",
	"mediumint": "int32",
	"integer":   "int32",
	"double":    "float32",
	"varchar":   "string",
	"char":      "string",
	"text":      "string",
	"longtext":  "string",
	"date":      "time.Time",
	"datetime":  "time.Time",
	//-----------------golang----------
	"float":   "float32",
	"float32": "float32",
	"float64": "float64",

	"uint":        "uint64",
	"uint8":       "uint8",
	"int8":        "int8",
	"int16":       "int16",
	"uint16":      "uint16",
	"int32":       "int32",
	"uint32":      "uint32",
	"int64":       "int64",
	"uint64":      "uint64",
	"string":      "string",
	"time":        "time.Time",
	"integers":    "int",
	"arraystring": "string",
	"timestamp":   "time.Time",
	"boolean":     "bool",
}
