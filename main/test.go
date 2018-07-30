package main

import (
	"os"
	"fmt"
)

func main() {
	f, err := os.Create("E:\\data222\\test.log")
	if err != nil {
		fmt.Println(err)
	}
	f.Close()
}
