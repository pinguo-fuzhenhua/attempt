package migration

import (
	"context"
	"log"
	"math"
	"video/migrationv2/db"
	"video/migrationv2/job"
	pjob "video/migrationv2/job"
	"video/migrationv2/svc"
)

type SyncJob interface {
	DealAndInsert(ctx context.Context, date any) error
	Read(ctx context.Context, page int64) (any, error)
	Count(ctx context.Context) (int64, error)
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
	count, err := job.Count(ctx)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < int(math.Ceil(float64(count)/float64(pjob.PageSize))); i++ {
		data, err := job.Read(ctx, int64(i))
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
}

func (mcs *SyncManager) StopSyncWorker() {
	if mcs.MainCancel != nil {
		mcs.MainCancel()
	}
}
