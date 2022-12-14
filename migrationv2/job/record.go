package job

import (
	"context"
	"io/fs"
	"io/ioutil"
	"log"
	"math"

	"video/migrationv2/db"
	"video/migrationv2/models"

	"github.com/pinguo-icc/go-base/v2/ierr"
	tapi "github.com/pinguo-icc/transaction-svc/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PageSize = 2000
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

func (rj *RecordJob) DealAndInsert(ctx context.Context, data []*models.BankAccountTransaction) error {
	for _, v := range data {
		ioutil.WriteFile("offset.log", []byte(v.ID.Hex()), fs.ModePerm)
		// fmt.Println("v.ID.Hex()", v.ID.Hex())
		switch v.Operation {
		case models.Sale:
			if _, err := rj.tranClient.Sale(ctx, &tapi.BankOperationRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(math.Abs(float64(v.Amount))),
				BankType:              string(v.BankType),
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			}); err != nil {
				if ie, ok := ierr.FromError(err); ok && ie.SubCode == 300004 {
					log.Println(ie.Reason)
					continue
				}
				return err
			}

		case models.Reload:
			if _, err := rj.tranClient.Reload(ctx, &tapi.BankOperationRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(math.Abs(float64(v.Amount))),
				BankType:              string(v.BankType),
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			}); err != nil {
				if ie, ok := ierr.FromError(err); ok && ie.SubCode == 300004 {
					log.Println(ie.Reason)
					continue
				}
				return err
			}
		case models.TransferOut:
			if v.Amount > 0 {
				continue
			}
			col := rj.WfDb.Collection(rj.LogColl)
			counterPart := &models.Log{}
			err := col.FindOne(ctx, models.M{"dataChannel": "juice", "opType": "purchasePreset", "transID": v.TransID, "amount": models.M{"$gt": 0}}).Decode(counterPart)
			if err != nil {
				return err
			}
			_, err = rj.tranClient.TransferOut(ctx, &tapi.TransferOutRequest{
				Scope:                 v.Scope,
				UserId:                v.UserID,
				DeviceId:              v.DeviceID,
				OriginalTransactionId: v.OriginTransactionID,
				Amount:                int32(math.Abs(float64(v.Amount))),
				BankType:              string(v.BankType),
				ToUserId:              counterPart.UserID,
				Reason:                v.Reason,
				Comment:               v.Comment,
				ForceCreatedAt:        v.CreatedAt.Unix(),
			})
			if err != nil {
				if ie, ok := ierr.FromError(err); ok && ie.SubCode == 300004 {
					log.Println(ie.Reason)
					continue
				}
				return err
			}
		}
	}
	return nil
}

func (wfj *RecordJob) Read(ctx context.Context, lastID models.ID) ([]*models.BankAccountTransaction, error) {
	col := wfj.WfDb.Collection(wfj.LogColl)
	logs := []*models.Log{}
	filter := models.M{
		"opType": models.M{"$ne": "initAccount"},
		"_id":    models.M{"$gt": lastID},
	}
	cur, err := col.Find(ctx, filter, options.Find().SetSort(bson.D{bson.E{Key: "_id", Value: 1}}).SetLimit(PageSize))
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &logs)
	if err != nil {
		return nil, err
	}
	res := []*models.BankAccountTransaction{}
	for _, log := range logs {
		transaction, err := log.ToTransaction()
		if err != nil {
			return nil, err
		}
		if transaction != nil {
			res = append(res, transaction)
		}
	}
	return res, err
}
