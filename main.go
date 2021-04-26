package main

import (
	"fmt"
	"github.com/1819997197/ws-tools/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println("cmd.Execute err:", err)
	}
}
