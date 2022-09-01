package job

import (
	"context"
	"fmt"
	"video/migration/db"
	"video/migration/models"

	tapi "github.com/pinguo-icc/transaction-svc/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordJob struct {
	*db.DB
	LogColl    string
	tranClient tapi.BankAccountClient
}

func NewRecordJob(db *db.DB, tran tapi.BankAccountClient) *RecordJob {
	return &RecordJob{
		DB:         db,
		LogColl:    "log",
		tranClient: tran,
	}
}

func (rj *RecordJob) DealAndInsert(ctx context.Context, data any) error {
	trans, _ := data.([]*models.BankAccountTransaction)
	for i, v := range trans {
		fmt.Println(i)
		switch v.Operation {
		case models.Sale:
			rj.tranClient.Sale(ctx, &tapi.BankOperationRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(v.Amount),
				BankType:              string(v.BankType),
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			})
		case models.Reload:
			rj.tranClient.Reload(ctx, &tapi.BankOperationRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(v.Amount),
				BankType:              string(v.BankType),
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			})
		case models.TransferOut:
			if v.Amount < 0 {
				continue
			}
			col := rj.WfDb.Collection(rj.LogColl)
			counterPart := &models.Log{}
			err := col.FindOne(ctx, models.M{"transID": v.TransID}).Decode(counterPart)
			if err != nil {
				return err
			}
			rj.tranClient.TransferOut(ctx, &tapi.TransferOutRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(v.Amount),
				BankType:              string(v.BankType),
				ToUserId:              counterPart.UserID,
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			})
		}
	}
	return nil
}

func (wfj *RecordJob) Read(ctx context.Context) (any, error) {
	col := wfj.WfDb.Collection(wfj.LogColl)
	logs := []*models.Log{}
	filter := models.M{
		"opType": models.M{
			"$nin": []string{"initAccount"},
		},
	}
	cur, err := col.Find(ctx, filter, options.Find().SetSort(bson.D{bson.E{Key: "_id", Value: 1}}).SetLimit(5000))
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &logs)
	if err != nil {
		return nil, err
	}
	res := []*models.BankAccountTransaction{}
	for _, log := range logs {
		transaction := log.ToTransaction()
		if transaction != nil {
			res = append(res, transaction)
		}
	}
	return res, err
}
