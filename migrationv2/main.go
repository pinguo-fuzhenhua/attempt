package main

import (
	"context"
	"flag"
	"log"
	"os"
	"video/migrationv2/db"
	"video/migrationv2/migration"
	"video/migrationv2/svc"
)

func main() {
	copt := db.InitConnOpt(flag.NewFlagSet(os.Args[0], flag.ExitOnError))
	db.InitConnection(copt)
	db := db.NewDb("waterfalls", "transaction-svc")
	clientSet, cancel := svc.NewClientSet(copt)
	sm := migration.NewSyncManager(db, clientSet)
	ctx, close := context.WithCancel(context.Background())
	sm.MainCancel = cancel
	defer func() {
		sm.StopSyncWorker()
		close()
		log.Println("---------------shutdown--------------")
	}()

	sm.RunSyncWorker(ctx)
}
