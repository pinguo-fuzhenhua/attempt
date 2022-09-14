package migration

import (
	"context"
	"fmt"
	"log"
	"math"
	"video/migrationv2/db"
	"video/migrationv2/job"
	pjob "video/migrationv2/job"
	"video/migrationv2/models"
	"video/migrationv2/svc"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncJob interface {
	DealAndInsert(ctx context.Context, date any) error
	Read(ctx context.Context, page int64, min, max models.ID) (any, error)
	Count(ctx context.Context, min, max models.ID) (int64, error)
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
	ids := []models.ID{primitive.NilObjectID, primitive.NewObjectID()}
	count, err := job.Count(ctx, ids[len(ids)-2], ids[len(ids)-1])
	if err != nil {
		log.Println(err)
		panic(err)
	}

	for count != 0 {
		for i := 0; i < int(math.Ceil(float64(count)/float64(pjob.PageSize))); i++ {
			fmt.Println("*******************************************************************", i, i, i)
			data, err := job.Read(ctx, int64(i), ids[len(ids)-2], ids[len(ids)-1])
			if err != nil {
				log.Println(err.Error())
				return
			}
			err = job.DealAndInsert(ctx, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
		ids = append(ids, primitive.NewObjectID())
		count, err = job.Count(ctx, ids[len(ids)-2], ids[len(ids)-1])
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
}
