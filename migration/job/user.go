package job

import (
	"context"
	"video/migration/db"
	"video/migration/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserJob struct {
	*db.DB
	UserColl string

	UserSnapshotColl string
	AccountColl      string
	IsMigrate        map[string]struct{}
}

func NewUserJob(db *db.DB) *UserJob {
	return &UserJob{
		DB:               db,
		UserColl:         "videobeats_user",
		UserSnapshotColl: "videobeats_user_snapshot",
		AccountColl:      "account",
		IsMigrate:        make(map[string]struct{}),
	}
}

// func (uj *UserJob) DealAndInsert(ctx context.Context, data any) error {
// 	col := uj.TranDb.Collection(uj.UserColl)
// 	if data, ok := data.([]interface{}); ok {
// 		return uj.TrTran.Do(ctx, func(ctx context.Context) error {
// 			_, err := col.InsertMany(ctx, data)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		})
// 	}
// 	return nil
// }

func (uj *UserJob) DealAndInsert(ctx context.Context, data any) error {
	col := uj.TranDb.Collection(uj.UserColl)
	if data, ok := data.([]interface{}); ok {
		_, err := col.InsertMany(ctx, data)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (uj *UserJob) Read(ctx context.Context) (any, error) {
	col := uj.WfDb.Collection(uj.AccountColl)
	accounts := []*models.Account{}
	cur, err := col.Find(ctx, primitive.M{})
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, err
	}
	res := []interface{}{}
	for _, account := range accounts {
		if _, ok := uj.IsMigrate[account.UserID]; !ok {
			user := account.ToUser()
			if user == nil {
				continue
			}
			res = append(res, user)
			uj.IsMigrate[account.UserID] = struct{}{}
		}
	}
	return res, err
}
