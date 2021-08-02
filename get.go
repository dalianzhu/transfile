package main

import (
	"context"
	"log"
	"os"

	"github.com/dalianzhu/transfile/proto"
	"google.golang.org/grpc"
)

func Get(address, code, filepath string) error {
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

	fi, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fi.Close()

	totalData := 0
	for {
		rsp, err := client.Recv()
		if err != nil {
			return err
		}
		if rsp.Head["op"] == "continue" {
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
