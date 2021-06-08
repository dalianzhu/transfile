package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	cmd := flag.String("m", "", "模式：可以为 agent/put/get")
	code := flag.String("c", "", "组编号")
	address := flag.String("a", "", "服务地址：127.0.0.1:9886")
	file := flag.String("f", "", "文件路径")
	flag.Parse()

	if os.Getenv("address") != "" {
		*address = os.Getenv("address")
	}

	if os.Getenv("code") != "" {
		*code = os.Getenv("code")
	}

	switch *cmd {
	case "agent":
		runAgent()
	case "get":
		err := Get(*address, *code, *file)
		if err != nil {
			log.Printf("get error:%v", err)
		}
	case "put":
		err := Put(*address, *code, *file)
		if err != nil {
			log.Printf("put error:%v", err)
		}
	}
}
