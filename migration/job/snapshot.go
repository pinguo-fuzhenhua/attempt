package job

import (
	"context"
	"video/migration/db"
	"video/migration/models"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapshotJob struct {
	*db.DB
	AccountColl string

	UserSnapshotColl string
}

func NewSnapshotJob(db *db.DB) *UserJob {
	return &UserJob{
		DB:               db,
		UserSnapshotColl: "videobeats_user_snapshot",
		AccountColl:      "account",
	}
}

// func (sj *SnapshotJob) DealAndInsert(ctx context.Context, data any) error {
// 	col := sj.TranDb.Collection(sj.UserSnapshotColl)
// 	if data, ok := data.([]interface{}); ok {
// 		return sj.TrTran.Do(ctx, func(ctx context.Context) error {
// 			_, err := col.InsertMany(ctx, data, options.InsertMany())
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		})
// 	}
// 	return nil
// }

func (sj *SnapshotJob) DealAndInsert(ctx context.Context, data any) error {
	col := sj.TranDb.Collection(sj.UserSnapshotColl)
	if data, ok := data.([]interface{}); ok {
		_, err := col.InsertMany(ctx, data, options.InsertMany())
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (sj *SnapshotJob) Read(ctx context.Context) (any, error) {
	col := sj.WfDb.Collection(sj.AccountColl)
	a := models.A{
		models.M{"$match": models.M{"deleted": 0}},
		models.M{"$group": models.M{
			"_id":        "$userID",
			"snapShotID": models.M{"$first": "$_id"},
			"account":    models.M{"$push": "$dataChannel"},
			"amount":     models.M{"$push": "$amount"},
			"createdAt":  models.M{"$first": "$createdAt"},
		},
		},
	}
	cur, err := col.Aggregate(ctx, a)
	if err != nil {
		return nil, err
	}
	mUsers := []*models.MiddleUser{}
	err = cur.All(ctx, &mUsers)
	if err != nil {
		return nil, err
	}
	res := []interface{}{}
	for _, account := range mUsers {
		userSnapshot := account.ToSnapshot()
		if userSnapshot == nil {
			continue
		}
		res = append(res, userSnapshot)
	}
	return res, err
}
