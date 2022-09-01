package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once                               sync.Once
	waterfallClient, transactionClient *mongo.Client
)

func getWaterfallClient(cfg *ConnCfg) {
	if waterfallClient == nil || transactionClient == nil {
		once.Do(func() {
			var err error
			ctx, _ := context.WithTimeout(context.Background(), time.Duration(cfg.TimeOut)*time.Second)
			waterfallClient, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.WaterFallDsn))
			if err != nil {
				log.Fatalf(err.Error())
			}

			transactionClient, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.TransactionDsn))
			if err != nil {
				log.Fatalf(err.Error())
			}
		})
	}
}

type SyncJob interface {
	Insert(*mongo.Client) error
	Read(*mongo.Client) error
}

func NewWaterFallJob() SyncJob {
	return &WaterFallJob{
		WaterFallDbName: "",
		WaterFallColl:   "",
	}
}

type WaterFallJob struct {
	WaterFallDbName string
	WaterFallColl   string
}

func (wfj *WaterFallJob) Insert(cli *mongo.Client) error {
	return nil
}

func (wfj *WaterFallJob) Read(cli *mongo.Client) error {
	return nil
}

func NewTransactionJob() SyncJob {
	return &TransactionJob{
		TransName:     "",
		WaterFallColl: "",
	}
}

type TransactionJob struct {
	TransName     string
	WaterFallColl string
}

func (tj *TransactionJob) Insert(cli *mongo.Client) error {
	return nil
}

func (tj *TransactionJob) Read(cli *mongo.Client) error {
	return nil
}

type ConnCfg struct {
	WaterFallDsn   string
	TransactionDsn string
	TimeOut        int64
}

func (c *ConnCfg) addFlag(f *flag.FlagSet) {
	f.StringVar(&c.WaterFallDsn, "waterfallDsn", "", "waterfall mongo dsn")
	f.StringVar(&c.TransactionDsn, "transactionDsn", "", "transaction mongo dsn")
	f.Int64Var(&c.TimeOut, "timeout", 2, "connecting time limit")
}

func (c *ConnCfg) validate() error {
	if c.WaterFallDsn == "" || c.TransactionDsn == "" {
		return errors.New("waterfall or transaction dns is nil")
	}
	return nil
}

func initConnOpt(f *flag.FlagSet) *ConnCfg {
	connOpt := new(ConnCfg)
	connOpt.addFlag(f)
	err := f.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf(err.Error())
	}
	return connOpt
}

type SyncManager struct {
	WaterfallClient *mongo.Client
	TransClient     *mongo.Client
	StopFunc        []context.CancelFunc
	SyncJobs        []SyncJob
	MainCancel      context.CancelFunc
	mux             sync.Mutex
}

type SyncOpt func(*SyncManager)

func WithSyncJobs(jobs ...SyncJob) SyncOpt {
	return func(sm *SyncManager) {
		sm.SyncJobs = jobs
	}
}

func NewSyncManager() *SyncManager {
	opts := []SyncOpt{
		WithSyncJobs(
			NewWaterFallJob(),
			NewTransactionJob(),
		),
	}
	syncManager := &SyncManager{
		WaterfallClient: waterfallClient,
		TransClient:     transactionClient,
		StopFunc:        []context.CancelFunc{},
		SyncJobs:        []SyncJob{},
	}
	for _, opt := range opts {
		opt(syncManager)
	}
	return syncManager
}

func (sm *SyncManager) RunOneJob(job SyncJob, wg *sync.WaitGroup) {
	// job.Read(sm.c)
	wg.Done()
}

func (mcs *SyncManager) StopSyncWorker() {
	for v := range mcs.StopFunc {
		mcs.StopFunc[v]()
	}
	if mcs.MainCancel != nil {
		mcs.MainCancel()
	}
}

type Log struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TS             time.Time          `bson:"ts"`                 //时间戳
	AppName        string             `bson:"appName"`            //产品名称
	DataChannel    string             `bson:"dataChannel"`        //数据渠道
	UserID         string             `bson:"userID,omitempty"`   //用户id
	DeviceID       string             `bson:"deviceID,omitempty"` //设备id
	OpType         string             `bson:"opType"`             //操作类型
	OriginalAmount int64              `bson:"originalAmount"`     //操作前数量
	Amount         int64              `bson:"amount"`             //操作数量
	Comment        string             `bson:"comment,omitempty"`  //描述
	TransID        string             `bson:"transID"`            //保证幂等的ID
}

type Account struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AppName     string             `bson:"appName"`
	DataChannel string             `bson:"dataChannel"`
	UserID      string             `bson:"userID"`
	DeviceID    string             `bson:"deviceID"`
	Amount      int64              `bson:"amount"`
	CreatedAt   time.Time          `bson:"createdAt"`
	BoundAt     time.Time          `bson:"boundAt"`
	Deleted     int                `bson:"deleted"`
}

type Transaction struct {
	ID                  primitive.ObjectID `bson:"_id"`
	TransactionID       string             `bson:"transactionId"`
	OriginTransactionID string             `bson:"originTransactionId"`
	Scope               string             `bson:"scope"`
	UserID              string             `bson:"userId"`
	Comment             string             `bson:"comment"`
	DeviceID            string             `bson:"deviceId"`
	Type                string             `bson:"type"`
	IsMockUser          bool               `bson:"isMockUser"`
	IsDeleted           bool               `bson:"isDeleted"`
	CreatedAt           time.Time          `bson:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt"`
}

type BankAccountTransactionSource struct {
	TransactionID string `bson:"transactionId"`
	Amount        int64  `bson:"amount"`
}

type CombinationChildren struct {
	BankType      BankType `bson:"bankType"`
	TransactionID string   `bson:"transactionId"`
	Amount        int64    `bson:"amount"`
	RestAmount    int64    `bson:"restAmount"`
}

type BankAccountTransaction struct {
	*Transaction         `bson:",inline"`
	BankType             BankType                                `bson:"bankType"` // 银行帐户类型,juice,points等
	Reason               string                                  `bson:"reason"`
	Amount               int64                                   `bson:"amount"`
	RestAmount           int64                                   `bson:"restAmount"`           // 可使用的数量, 默认restAmount=amount, 对于reload,capture,refund,transferIn等可用的钱, 需要记录已经使用的数量
	SourceTransactionIDs map[string]BankAccountTransactionSource `bson:"sourceTransactionIds"` // 需要可追述, 记录该次交易的上级
	ChildTransactionIDs  map[string]BankAccountTransactionSource `bson:"childTransactionIds"`  // 记录钱的去处
	Operation            BankAccountOperation                    `bson:"operation"`
	CombinationChildren  []*CombinationChildren                  `bson:"combinationChildren,omitempty"` // 组合账户时, 用于子帐号信息保存
}

func main() {
	copt := initConnOpt(flag.NewFlagSet(os.Args[0], flag.ExitOnError))
	getWaterfallClient(copt)
	sm := NewSyncManager()
	ctxMain, cancel := context.WithTimeout(context.Background(), time.Second*10)
	sm.MainCancel = cancel
	for i := range sm.SyncJobs {
		job := sm.SyncJobs[i]
		sm.mux.Lock()

		go sm.RunOneJob(job, wg)
	}
	wg.Wait()
}
