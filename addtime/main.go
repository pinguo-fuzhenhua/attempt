package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbname     = "idphoto"
	collName   = "origin"
	BatchLimit = 2000
	offset     = 20
)

type option struct {
	MongoDNS string
	Timeout  int
}

func (o *option) validate() error {
	if o.MongoDNS == "" {
		return errors.New("please set mongo connect uri")
	}

	return nil
}

func (o *option) addFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.MongoDNS, "idPhotoDsn", "", "the mongoDB connect address")
	fs.IntVar(&o.Timeout, "timeout", 2, "the exec timeout setting,default 1 minute")
}

func initOptions(fs *flag.FlagSet, args ...string) (*option, error) {
	o := new(option)
	o.addFlags(fs)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o, nil
}

type DocOrigin struct {
	ID        primitive.ObjectID `bson:"_id"`       // 源图ID
	DeviceID  string             `bson:"device_id"` // 设备ID
	UserID    string             `bson:"user_id"`   // 用户ID
	Key       string             `bson:"key"`       // 七牛源图key
	CreatedAt time.Time          `bson:"createdAt"` // 创建时间
}

func modify(cli *mongo.Client, i int64, mu sync.Mutex) {
	coll := cli.Database(dbname).Collection(collName)
	// lastID := primitive.NilObjectID
	// if body, err := os.ReadFile("offset.log"); err == nil {
	// 	lastID, err = primitive.ObjectIDFromHex(string(body))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// coll.UpdateMany(context.Background(), primitive.M{})
	fmt.Printf("第%d个线程启动\n", i)
	cur, err := coll.Find(context.Background(), primitive.M{}, options.Find().SetSort(bson.D{bson.E{Key: "_id", Value: 1}}).SetSkip(i*int64(BatchLimit)).SetLimit(int64(BatchLimit)))
	if err != nil {
		panic(err)
	}

	photos := []*DocOrigin{}
	err = cur.All(context.Background(), &photos)
	if err != nil {
		panic(err)
	}

	for _, v := range photos {
		v.CreatedAt = v.ID.Timestamp()
		_, err := coll.UpdateByID(context.Background(), v.ID, bson.D{bson.E{Key: "$set", Value: v}})
		if err != nil {
			panic(err)
		}
	}
	if len(photos) == 0 {
		return
	}
	mu.Lock()
	ioutil.WriteFile("offset.log", []byte(photos[len(photos)-1].ID.Hex()), fs.ModePerm)
	mu.Unlock()
}

func main() {
	o, err := initOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(o.Timeout)*time.Minute)
	defer cancel()

	db, err := mongo.Connect(ctx, options.Client().ApplyURI(o.MongoDNS))
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("start modify")
	coll := db.Database(dbname).Collection(collName)
	// ret, err := coll.UpdateMany(context.Background(), primitive.M{"createAt": primitive.M{"$exists": true}}, primitive.M{"$unset": primitive.M{"createAt": nil}})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(ret.UpsertedCount)
	r := DocOrigin{}
	_ = coll.FindOne(context.Background(), primitive.M{}, options.FindOne().SetSkip(int64(offset)).SetSort(bson.D{bson.E{Key: "_id", Value: -1}})).Decode(&r)
	fmt.Println(r)
	// cur, err := coll.CountDocuments(context.Background(), primitive.M{"createdAt": primitive.M{"$exists": false}})
	// fmt.Println(cur)
	// cur.Decode(&r)
	// fmt.Println(r)
	// var mu sync.Mutex
	// wg := sync.WaitGroup{}
	// for i := 0; i < 121; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		if i != 0 {
	// 			time.Sleep(time.Duration(i) * time.Second)
	// 		}
	// 		modify(db, int64(i), mu)
	// 		fmt.Printf("第%d个线程同步数据同步已完成\n", i)
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()
	// log.Println("modify end")
}
