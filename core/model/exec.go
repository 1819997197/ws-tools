package model

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

func GenerationModel(distPath, packageName, connect, destTableName string) error {
	// 1.获取表名
	var err error
	var destTable []string
	if destTableName != "" {
		tableNames := strings.Split(destTableName, ",")
		for _, t := range tableNames {
			t = strings.TrimSpace(t)
			if t != "" {
				destTable = StringArrayAppend(destTable, t)
			}
		}
	}

	// 2.连接数据库
	var db *sqlx.DB
	db, err = sqlx.Connect("mysql", connect)
	if err != nil {
		log.Print("connect db err:", err)
		return errors.New("connect db err")
	}

	// 3.获取表结构信息
	var modelInit = &ModelInit{
		State:   []string{"State"},
		Created: []string{"Created"},
		Updated: []string{"Updated"},
		Deleted: []string{"Deleted", "Deprecated"},
	}
	tableList, err := Tables(db, destTable, modelInit)
	if err != nil {
		log.Print("get table infos err:", err)
		return errors.New("get table infos err")
	}

	// 4.遍历生成model层代码
	for _, t := range tableList {
		doc := TableDoc(t, packageName)
		header := Header(t, packageName)
		fileName := distPath + "/" + t.Table + ".go"
		err = FileWrite(fileName, header+doc)
		if err != nil {
			log.Println("FileWrite "+fileName+" 发生错误: ", err.Error())
		}
	}

	return nil
}
