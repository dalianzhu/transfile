package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/dalianzhu/transfile/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func runAgent() {
	var keepAliveArgs = keepalive.ServerParameters{
		Time:              1 * time.Minute,
		Timeout:           15 * time.Second,
		MaxConnectionIdle: 5 * time.Minute,
	}
	port := 9886
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
		return
	}
	logrus.Infof("runAgent listen:%v", port)

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepAliveArgs),
		grpc.MaxSendMsgSize(1024*1024*4), // 最大消息4M
		// grpc.ReadBufferSize(1024),
		// grpc.WriteBufferSize(1024),
		// grpc.MaxConcurrentStreams(10),
	)

	svc := &TransFile{}

	proto.RegisterTranFileAgentServer(s, svc)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}

type TransFile struct {
}

var insMap sync.Map

type Ins struct {
	lk     *sync.Mutex
	Code   int
	GetSvr proto.TranFileAgent_GetServer
}

var lk sync.Mutex

func (t *TransFile) Put(svr proto.TranFileAgent_PutServer) (err error) {
loop:
	for {
		data, err := svr.Recv()
		if err != nil {
			log.Printf("put recv error:%v", err)
			break loop
		}
		op := data.Head["op"]
		code := data.Head["code"]
		blk := data.Head["blk"]
		end := data.Head["end"]
		log.Printf("op:%v code:%v blk:%v len:%v", op, code, blk, len(data.Data))
		v, ok := insMap.Load(code)
		if !ok {
			ins := new(Ins)
			ins.lk = new(sync.Mutex)
			log.Printf("code:%v create", code)
			insMap.Store(code, ins)
			defer func() {
				log.Printf("code:%v exit", code)
				insMap.Delete(code)
			}()
			svr.Send(&proto.Data{
				Head: map[string]string{
					"op": "wait",
				},
			})
		} else {
			ins := v.(*Ins)
			if ins.GetSvr != nil {
				// 给get的客户端发送文件
				getOp := "continue"
				if end == "end" {
					getOp = "end"
				}
				err := ins.GetSvr.Send(&proto.Data{
					Head: map[string]string{
						"op":  getOp,
						"blk": blk,
					},
					Data: data.Data,
				})
				if err != nil {
					log.Printf("put getsvr send error:%v", err)
					break loop
				}
				// 告诉put的客户端，继续发送文件
				// 如果put发了end过来，就回一个end过去
				svr.Send(&proto.Data{
					Head: map[string]string{
						"op": getOp,
					},
				})
			} else {
				svr.Send(&proto.Data{
					Head: map[string]string{
						"op": "wait",
					},
				})
			}
		}
	}
	return nil
}

func (t *TransFile) Get(svr proto.TranFileAgent_GetServer) (err error) {
	for {
		data, err := svr.Recv()
		if err != nil {
			return err
		}
		code := data.Head["code"]
		v, ok := insMap.Load(code)
		if !ok {
			return fmt.Errorf("%v is not exist", code)
		}
		ins := v.(*Ins)
		if ins.GetSvr == nil {
			ins.GetSvr = svr
		}
	}
}
