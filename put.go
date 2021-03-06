package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dalianzhu/transfile/proto"
	"google.golang.org/grpc"
)

func Put(address, code, filePath string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := proto.NewTranFileAgentClient(conn)
	client, err := c.Put(context.Background())
	if err != nil {
		return err
	}

	currentBlk := 0
	fi, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	buf := make([]byte, 1024*1024) //一次读取多少个字节
	willEnd := false
	totalSend := 0
loop:
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		var chunk []byte
		if n == 0 {
			willEnd = true
			chunk = nil
		} else {
			chunk = buf[:n]
		}
	send:
		for {
			currentBlk++
			zipChuck, err := GZipData(chunk)
			if err != nil {
				return err
			}
			data := &proto.Data{
				Head: map[string]string{
					"op":       "put",
					"code":     code,
					"blk":      fmt.Sprint(currentBlk),
					"end":      "",
					"filePath": filePath,
				},
				Data: zipChuck,
			}

			if willEnd {
				data.Head["end"] = "end"
			}
			err = client.Send(data)
			if err != nil {
				return err
			}

			rsp, err := client.Recv()
			if err != nil {
				return err
			}
			rspOp := rsp.Head["op"]
			switch rspOp {
			case "wait":
				time.Sleep(time.Second)
				continue send
			case "continue":
				totalSend += len(chunk)
				log.Printf("send: %v\n", totalSend)
				break send
			case "end":
				client.CloseSend()
				break loop
			}
		}
	}
	return nil
}
