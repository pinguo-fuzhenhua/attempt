package db

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once                               sync.Once
	WaterfallClient, TransactionClient *mongo.Client
)

func InitConnection(cfg *ConnCfg) {
	if WaterfallClient == nil || TransactionClient == nil {
		once.Do(func() {
			var err error
			ctx, _ := context.WithTimeout(context.Background(), time.Duration(cfg.TimeOut)*time.Second)
			WaterfallClient, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.WaterFallDsn))
			if err != nil {
				log.Fatalf(err.Error())
				panic(err)
			}

			TransactionClient, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.TransactionDsn))
			if err != nil {
				log.Fatalf(err.Error())
				panic(err)
			}
		})
	}
}

type ConnCfg struct {
	WaterFallDsn       string
	TransactionDsn     string
	TransactionSvcAddr string
	TimeOut            int64
}

func (c *ConnCfg) addFlag(f *flag.FlagSet) {
	f.StringVar(&c.WaterFallDsn, "waterfallDsn", "", "waterfall mongo dsn")
	f.StringVar(&c.TransactionDsn, "transactionDsn", "", "transaction mongo dsn")
	f.StringVar(&c.TransactionSvcAddr, "transactionSvcAddr", "", "transactionSvcAddr")
	f.Int64Var(&c.TimeOut, "timeout", 2, "connecting time limit")
}

func (c *ConnCfg) validate() error {
	if c.WaterFallDsn == "" || c.TransactionDsn == "" {
		return errors.New("waterfall or transaction dns is nil")
	}
	return nil
}

func InitConnOpt(f *flag.FlagSet) *ConnCfg {
	connOpt := new(ConnCfg)
	connOpt.addFlag(f)
	err := f.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = connOpt.validate()
	if err != nil {
		log.Fatal(err.Error())
	}
	return connOpt
}

type DB struct {
	WfDb, TranDb *mongo.Database
}

func NewDb(wfDbName, tranDbName string) *DB {
	wfDb := WaterfallClient.Database(wfDbName)
	tranDb := TransactionClient.Database(tranDbName)
	return &DB{
		WfDb:   wfDb,
		TranDb: tranDb,
	}
}
