package svc

import (
	"context"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"

	"video/migrationv2/db"

	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	kierr "github.com/pinguo-icc/kratos-library/v2/ierr"
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
		kgrpc.WithTimeout(60*time.Second),
		kgrpc.WithOptions(
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:    time.Second * 3,
					Timeout: time.Second * 3,
				},
			),
		),
		kgrpc.WithMiddleware(
			kierr.GRPCClientMiddleware(),
		),
	)
	if err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}

	return conn, nil
}
