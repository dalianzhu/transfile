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

	for {
		rsp, err := client.Recv()
		if err != nil {
			return err
		}
		log.Printf("revc blk:%v", rsp.Head["blk"])
		if rsp.Head["op"] == "continue" {
			fi.Write(rsp.Data)
			continue
		} else if rsp.Head["op"] == "end" {
			log.Printf("revc end")
			return nil
		}
	}
}
