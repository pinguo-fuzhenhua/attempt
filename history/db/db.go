package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"video/history/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	D        = bson.D
	E        = bson.E
	M        = bson.M
	A        = bson.A
	ObjectID = primitive.ObjectID
)

type dbEntity struct {
	dbName                 string
	operateHistoryCollName string
	operateEntityCollName  string
}

var Executor dbEntity

func init() {
	Executor = dbEntity{
		dbName:                 "operate-history",
		operateHistoryCollName: "operateHistory",
		operateEntityCollName:  "testEntityOperate",
	}
}

func (de *dbEntity) Run(ctx context.Context, client *mongo.Client) {
	log.Printf("======run on db ===========\n")
	db := client.Database(de.dbName)
	ohColl := db.Collection(de.operateHistoryCollName)
	eoColl := db.Collection(de.operateEntityCollName)
	wg := &sync.WaitGroup{}
	var ohs []*models.OperateHistory
	var next int64
	count := 1
	for count > 0 {
		cur, err := ohColl.Find(ctx, M{},
			options.Find().
				SetSort(bson.D{bson.E{Key: "_id", Value: 1}}).
				SetSkip(int64(next*500)).
				SetLimit(500),
		)
		if err != nil {
			log.Println("method=historySplite,", fmt.Sprintf("err=%s", err.Error()))
		}
		err = cur.All(ctx, &ohs)
		if err != nil {
			log.Println("method=historySplite,", fmt.Sprintf("err=%s", err.Error()))
		}
		count = len(ohs)
		if count > 0 {
			wg.Add(1)
			go func(ohs []*models.OperateHistory) {
				defer wg.Done()
				eos := []any{}
				for _, v := range ohs {
					eos = append(eos, &models.EntityOperate{
						ID:          primitive.NewObjectIDFromTimestamp(v.OperateTime),
						OperateTime: v.OperateTime,
						EntityType:  v.EntityType,
						EntityName:  v.EntityName,
						EntityID:    v.EntityID,
						TraceID:     v.TraceID,
						Scope:       v.Scope,
						Environment: v.Environment,
						OperateType: v.OperateType,
					})
				}
				_, err := eoColl.InsertMany(context.Background(), eos)
				if err != nil {
					log.Println("method=historySplite, func=entityOperate", fmt.Sprintf("err=%s", err.Error()))
				}
			}(ohs)
		}
		next += 1
	}
	wg.Wait()

	// 更新主表，去掉废弃字段
	d := M{
		"$unset": M{"entityType": "", "entityName": "", "entityID": ""},
	}
	_, err := ohColl.UpdateMany(ctx, M{}, d)
	if err != nil {
		log.Println("method=historySplite, func=entityOperate, remove abandoned fields from operate_history", fmt.Sprintf("err=%s", err.Error()))
	}
	log.Printf("======run finished =========== \n\n")
}
