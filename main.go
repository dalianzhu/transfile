package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	code := flag.String("c", "", "组编号")
	address := flag.String("a", "", "服务地址：127.0.0.1:9886")

	flag.Usage = func() {
		fmt.Printf(`Usage of %s:
transfile [cmd: put/get/agent] [filepath:./hello.tgz] -a address -c groupCode
收发文件需要先建立服务端
transfile -a 127.0.0.1:9886 agent

发送文件：
transfile -a 127.0.0.1:9886 -c 1 put hello.tgz 

接收文件：
transfile -a 127.0.0.1:9886 -c 1 get hello.tgz 

环境变量（可选）
export f_address=127.0.0.1:9886
export f_code=1

参数：
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	cmd := flag.Arg(0)
	filePath := flag.Arg(1)

	if os.Getenv("f_address") != "" && *address == "" {
		*address = os.Getenv("f_address")
	}
	if os.Getenv("f_code") != "" && *code == "" {
		*code = os.Getenv("f_code")
	}

	fmt.Printf("address:%v, code:%v, cmd:%v, file:%v\n", *address, *code, cmd, filePath)
	if *address == "" {
		fmt.Println("地址为空")
		return
	}

	switch cmd {
	case "agent":
		runAgent()
	case "get":
		if *code == "" {
			fmt.Println("组编号为空")
			return
		}
		err := Get(*address, *code, filePath)
		if err != nil {
			log.Printf("get error:%v", err)
		}
	case "put":
		if err := checkFileAndCode(*code, filePath); err != nil {
			fmt.Println(err)
			return
		}
		err := Put(*address, *code, filePath)
		if err != nil {
			log.Printf("put error:%v", err)
		}
	default:
		fmt.Println("命令错误")
	}
}

func checkFileAndCode(code, file string) error {
	if code == "" {
		return fmt.Errorf("组编号为空")
	}
	if file == "" {
		return fmt.Errorf("发送文件地址为空")
	}
	return nil
}
