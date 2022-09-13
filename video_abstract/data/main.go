package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbEntity struct {
	dbName string
	coll   []string
	scope  map[string][]string
}

func (e *dbEntity) dbNames() []string {
	names := []string{}
	for k, v := range e.scope {
		for _, env := range v {
			names = append(names, fmt.Sprintf("%s_%s_%s", k, env, e.dbName))
		}
	}
	return names
}

func (e *dbEntity) run(ctx context.Context, client *mongo.Client) {
	for _, v := range e.dbNames() {
		log.Printf("======run on db %s ===========\n", v)

		db := client.Database(v)
		for _, c := range e.coll {
			col := db.Collection(c)
			_ = col
			if err := setSoftDeleteField(ctx, col, 1, "oppos_material"); err != nil {
				log.Println(err)
			}

			if err := setSoftDeleteField(ctx, col, 3, "oppos_activity"); err != nil {
				log.Println(err)
			}
		}

		log.Printf("======run finished =========== \n\n")
	}
}

func setSoftDeleteField(ctx context.Context, coll *mongo.Collection, sort int, typeName string) error {
	doc := primitive.M{
		"code":                "tag",
		"name":                "标签",
		"comment":             "",
		"required":            false,
		"isSystematized":      true,
		"isLocalized":         false,
		"sysFieldType":        9,
		"customizedFieldType": 10,
		"tagField": primitive.M{
			"tagType": typeName,
		},
	}
	filter := primitive.M{"fields.code": primitive.M{"$ne": "tag"}, "type": sort}
	update := primitive.M{"$push": primitive.M{"fields": doc}}

	u, err := coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	log.Printf("%s %d doc has been set tagF field.\n", coll.Name(), u.ModifiedCount)

	return nil
}

func findSomething(ctx context.Context, coll *mongo.Collection) error {
	d := primitive.M{}

	cur, err := coll.CountDocuments(ctx, d)
	if err != nil {
		return err
	}

	log.Print(cur)

	return nil
}

var execs []dbEntity

func init() {
	execs = []dbEntity{
		{
			dbName: "field-definitions",
			coll:   []string{"fields_definition"},
			scope: map[string][]string{
				"videobeats": {"prod", "operation", "dev", "qa"},
				"camera360":  {"prod", "operation", "dev", "qa"},
				"idphoto":    {"prod", "operation", "dev", "qa"},
				"mix":        {"prod", "operation", "dev", "qa"},
				"salad":      {"prod", "operation", "dev", "qa"},
				"inface":     {"prod", "operation", "dev", "qa"},
				"icc":        {"prod", "operation", "dev", "qa"}},
		},
	}
}

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
	fs.StringVar(&o.MongoDNS, "mongo_dns", "", "the mongoDB connect address")
	fs.IntVar(&o.Timeout, "timeout", 4, "the exec timeout setting,default 1 minute")
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

func exec(ctx context.Context, cli *mongo.Client) {
	for _, v := range execs {
		v.run(ctx, cli)
	}

}

type AA struct {
	A []string
	B map[string]int
	C string
}

func main() {
	a := os.Args
	_ = a
	o, err := initOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(o.Timeout)*time.Minute)
	defer cancel()

	dbCli, err := mongo.Connect(ctx, options.Client().ApplyURI(o.MongoDNS))
	if err != nil {
		log.Fatal(err)
	}

	exec(ctx, dbCli)
}
