package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/dalianzhu/transfile/proto"
	"google.golang.org/grpc"
)

func Get(address, code, filePath string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := proto.NewTranFileAgentClient(conn)
	client, err := c.Get(context.Background())
	if err != nil {
		return err
	}

	err = client.Send(&proto.Data{
		Head: map[string]string{
			"op":   "get",
			"code": code,
		},
	})
	if err != nil {
		return err
	}

	var fi *os.File

	totalData := 0
	for {
		rsp, err := client.Recv()
		if err != nil {
			return err
		}
		if rsp.Head["op"] == "continue" {
			// 初始化文件
			if fi == nil {
				if filePath == "" {
					filePath = rsp.Head["filePath"]
				}
				_, fileName := path.Split(filePath)
				fi, err = os.Create(fileName)
				if err != nil {
					return err
				}
				defer fi.Close()
			}
			data, err := GUnzipData(rsp.Data)
			if err != nil {
				return err
			}
			fi.Write(data)
			totalData += len(data)
			log.Printf("revc:%v", totalData)
			continue
		} else if rsp.Head["op"] == "end" {
			log.Printf("revc end")
			return nil
		}
	}
}
