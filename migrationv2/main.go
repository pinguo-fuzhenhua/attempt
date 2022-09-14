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
	db := db.NewDb("waterfalls")
	clientSet, close := svc.NewClientSet(copt)
	sm := migration.NewSyncManager(db, clientSet)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		close()
		log.Println("---------------shutdown--------------")
	}()

	sm.RunSyncWorker(ctx)
}
