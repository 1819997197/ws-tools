package cmd

import (
	"fmt"
	"github.com/1819997197/ws-tools/core/model"
	"github.com/spf13/cobra"
	"log"
)

var distPath, packageName, connect, destTableName string

var modelCmd = &cobra.Command{
	Use:   "sql",
	Short: "Table structure auto generation model",
	Long:  "Table structure auto generation model",
	Run: func(cmd *cobra.Command, args []string) {
		if connect == "" {
			fmt.Println("please enter db connect dsn")
			return
		}
		err := model.GenerationModel(distPath, packageName, connect, destTableName)
		if err != nil {
			log.Print("generate fail")
			return
		}
		log.Print("generate success")
	},
}

func init() {
	modelCmd.Flags().StringVarP(&distPath, "dist", "", "./models", "model层代码生产目录")
	modelCmd.Flags().StringVarP(&packageName, "pkg", "", "models", "生成的代码与src的相对路径")
	modelCmd.Flags().StringVarP(&connect, "conn", "", "", "数据库连接dsn user:pwd@tcp(ip:port)/table?charset=utf8&parseTime=true")
	modelCmd.Flags().StringVarP(&destTableName, "table", "", "", "所需生成的表，用逗号分割(默认导出所有的表)")
}
