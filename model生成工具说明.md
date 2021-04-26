# model代码生成工具使用说明

## 一、安装
```
1.开启go mod模式: 
$ export GO111MODULE=on GOPROXY=https://goproxy.cn

2.引入包
package main

import (
	_ "github.com/1819997197/ws-tools"
)

3.生成mod文件
$ go mod init tools

4.安装
$ go install github.com/1819997197/ws-tools
安装成功后，会在$GOBIN目录生成一个二进制文件(ws-tools)
```

## 二、工具使用

#### 1.查看使用帮助
```
$ ws-tools sql --help
Table structure auto generation model

Usage:
sql [flags]

Flags:
--conn string    数据库连接dsn user:pwd@tcp(ip:port)/table?charset=utf8&parseTime=true
--dist string    model层代码生产目录 (default "./models")
-h, --help           help for sql
--pkg string     生成的代码与src的相对路径 (default "models")
--table string   所需生成的表，用逗号分割(默认导出所有的表)
```

#### 2.生成model文件
```
$ ws-tools sql --conn="user:pwd@tcp(ip:port)/db?charset=utf8&parseTime=true"
// 默认会在执行命令的当前目录的models目录下生成model文件(models目录需要存在)
```