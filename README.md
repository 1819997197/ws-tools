# ws-tools 工具集

## model层代码生产工具
```
1.查看命令
$ go run main.go sql --help
Table structure auto generation model

Usage:
   sql [flags]

Flags:
      --conn string    数据库连接dsn user:pwd@tcp(ip:port)/table?charset=utf8&parseTime=true
      --dist string    model层代码生产目录 (default "./models")
  -h, --help           help for sql
      --pkg string     生成代码的包名 (default "models")
      --table string   所需生成的表，用逗号分割(默认导出所有的表)

2.生成代码
go run main.go sql --conn="user:pwd@tcp(ip:port)/db?charset=utf8&parseTime=true"

3.具体使用细节见<model生成工具说明.md>
```

注：获取表结构信息部分fork于https://github.com/laixyz/go-mysql-model-creator.git(修复了部分获取索引bug)
