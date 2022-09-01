package svc

import (
	"context"
	"log"
	"video/migration/db"

	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	tapi "github.com/pinguo-icc/transaction-svc/api"
	"google.golang.org/grpc"
)

type ClientSet struct {
	TransactionClient tapi.BankAccountClient
}

func NewClientSet(cfg *db.ConnCfg) (*ClientSet, func()) {
	tranConn, err := newConnection(cfg.TransactionSvcAddr)
	if err != nil {
		return nil, nil
	}
	cancel := func() {
		tranConn.Close()
	}

	cs := &ClientSet{
		TransactionClient: tapi.NewBankAccountClient(tranConn),
	}

	return cs, cancel
}

func newConnection(addr string) (*grpc.ClientConn, error) {
	conn, err := kgrpc.DialInsecure(
		context.TODO(),
		kgrpc.WithEndpoint(addr),
	)
	if err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}

	return conn, nil
}
