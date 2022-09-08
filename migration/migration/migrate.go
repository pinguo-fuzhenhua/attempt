package migration

import (
	"context"
	"log"
	"sync"
	"video/migration/db"
	"video/migration/job"
	"video/migration/svc"
)

const (
	TotalPage = 10
)

type SyncJob interface {
	DealAndInsert(ctx context.Context, date any) error
	Read(ctx context.Context, page int64) (any, error)
}

type SyncManager struct {
	StopFunc   []context.CancelFunc
	SyncJobs   []SyncJob
	MainCancel context.CancelFunc
	mux        sync.Mutex
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
			// job.NewUserJob(db),
			// job.NewLogJob(db),
			// job.NewSnapshotJob(db),
			job.NewRecordJob(db, clients.TransactionClient),
		),
	}
	syncManager := &SyncManager{
		StopFunc: []context.CancelFunc{},
		SyncJobs: []SyncJob{},
	}
	for _, opt := range opts {
		opt(syncManager)
	}
	return syncManager
}

func (sm *SyncManager) RunSyncWorker(ctx context.Context) {
	wg := &sync.WaitGroup{}
	log.Println("--------------------start sync-------------------")
	for i := range sm.SyncJobs {
		wg.Add(1)
		job := sm.SyncJobs[i]
		go func() {
			sm.RunOneJob(ctx, job)
			wg.Done()
		}()
	}
	wg.Wait()
	sm.MainCancel()
}

func (sm *SyncManager) RunOneJob(ctx context.Context, job SyncJob) {
	for i := 1; i <= TotalPage; i++ {
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
	for v := range mcs.StopFunc {
		mcs.StopFunc[v]()
	}
	if mcs.MainCancel != nil {
		mcs.MainCancel()
	}
}
