package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
	"video/migration/db"
	"video/migration/migration"
	"video/migration/svc"
)

func main() {
	copt := db.InitConnOpt(flag.NewFlagSet(os.Args[0], flag.ExitOnError))
	db.InitConnection(copt)
	db := db.NewDb("waterfalls", "transaction-svc")
	clientSet, cancel := svc.NewClientSet(copt)
	sm := migration.NewSyncManager(db, clientSet)
	ctx, close := context.WithTimeout(context.Background(), time.Second*1000)
	sm.MainCancel = cancel
	defer func() {
		sm.StopSyncWorker()
		close()
		log.Println("---------------shutdown--------------")
	}()

	sm.RunSyncWorker(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("---------------sync end--------------")
			return
		}
	}
}
