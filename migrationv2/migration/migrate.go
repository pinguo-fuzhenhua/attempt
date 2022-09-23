package migration

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"video/migrationv2/db"
	"video/migrationv2/job"
	"video/migrationv2/models"
	"video/migrationv2/svc"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncJob interface {
	DealAndInsert(ctx context.Context, date []*models.BankAccountTransaction) error
	Read(ctx context.Context, lastId models.ID) ([]*models.BankAccountTransaction, error)
}

type SyncManager struct {
	SyncJobs   []SyncJob
	MainCancel context.CancelFunc
}

type SyncOpt func(*SyncManager)

func WithSyncJobs(jobs ...SyncJob) SyncOpt {
	return func(sm *SyncManager) {
		sm.SyncJobs = jobs
	}
}

func NewSyncManager(db *db.DB, clients *svc.ClientSet) *SyncManager {
	opts := []SyncOpt{
		WithSyncJobs(
			job.NewRecordJob(db, clients.TransactionClient),
		),
	}
	syncManager := &SyncManager{}
	for _, opt := range opts {
		opt(syncManager)
	}
	return syncManager
}

func (sm *SyncManager) RunSyncWorker(ctx context.Context) {
	log.Println("--------------------start sync-------------------")
	job := sm.SyncJobs[0]
	sm.RunOneJob(ctx, job)
}

func (sm *SyncManager) RunOneJob(ctx context.Context, job SyncJob) {
	lastID := primitive.NilObjectID
	if body, err := os.ReadFile("offset.log"); err == nil {
		lastID, err = primitive.ObjectIDFromHex(string(body))
		if err != nil {
			panic(err)
		}
	}
	for {
		fmt.Printf("==============  lastID=%s =============", lastID.Hex())
		data, err := job.Read(ctx, lastID)
		if err != nil {
			log.Panic("read", err)
		}
		if len(data) == 0 {
			log.Printf("lastID=%s 没有新数据产生\r\n", lastID.Hex())
			time.Sleep(1 * time.Second) // 休眠1s等待新数据产生
			continue
		}
		lastID = data[len(data)-1].ID

		err = job.DealAndInsert(ctx, data)
		if err != nil {
			log.Panic("DealAndInsert", err)
		}
	}
}
