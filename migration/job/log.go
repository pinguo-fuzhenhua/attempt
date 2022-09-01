package job

import (
	"context"
	"video/migration/db"
	"video/migration/models"
)

const (
	MaxCapacity = 100
)

type WaterFallJob struct {
	*db.DB
	FreeCountColl      string
	LogColl            string
	JuiceGiftColl      string
	JuicePurchasedColl string
	JuiceColl          string
}

func NewLogJob(db *db.DB) *WaterFallJob {
	return &WaterFallJob{
		DB:                 db,
		LogColl:            "log",
		FreeCountColl:      "videobeats_bank_free_count_transaction",
		JuiceGiftColl:      "videobeats_bank_juice_gift_transaction",
		JuicePurchasedColl: "videobeats_bank_juice_purchased_transaction",
		JuiceColl:          "videobeats_bank_juice_transaction",
	}
}

type TranContainer struct {
	items  [MaxCapacity]interface{}
	length int
}

func (tc *TranContainer) reset() {
	tc.length = 0
	tc.items = [MaxCapacity]interface{}{}
}

func (wfj *WaterFallJob) CollClassify(ctx context.Context,
	collName string,
	data *models.BankAccountTransaction,
	trans *TranContainer,
) error {
	var err error
	trans.items[trans.length] = data
	trans.length++
	if trans.length == 100 {
		err = wfj.BatchInsert(ctx, collName, trans.items[:])
		if err != nil {
			return err
		}
		trans.reset()
	}
	return err
}

func (wfj *WaterFallJob) BatchInsert(ctx context.Context, collName string, data []interface{}) error {
	_, err := wfj.TranDb.Collection(collName).InsertMany(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// func (wfj *WaterFallJob) DealAndInsert(ctx context.Context, data any) error {
// 	return wfj.TrTran.Do(ctx, func(ctx context.Context) error {
// 		var err error
// 		freeCountTran := new(TranContainer)
// 		juiceTran := new(TranContainer)
// 		juiceGiftTran := new(TranContainer)
// 		juicePurchasedTran := new(TranContainer)
// 		if data, ok := data.([]*models.BankAccountTransaction); ok {
// 			for i := range data {
// 				switch data[i].BankType {
// 				case models.FreeCount:
// 					err = wfj.CollClassify(ctx, wfj.FreeCountColl, data[i], freeCountTran)
// 				case models.Juice:
// 					err = wfj.CollClassify(ctx, wfj.JuiceColl, data[i], juiceTran)
// 				case models.JuiceGift:
// 					err = wfj.CollClassify(ctx, wfj.JuiceGiftColl, data[i], juiceGiftTran)
// 				case models.JuicePurchased:
// 					err = wfj.CollClassify(ctx, wfj.JuicePurchasedColl, data[i], juicePurchasedTran)
// 				}
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		if freeCountTran.length > 0 {
// 			err = wfj.BatchInsert(ctx, wfj.FreeCountColl, freeCountTran.items[:])
// 		} else if juiceTran.length > 0 {
// 			err = wfj.BatchInsert(ctx, wfj.JuiceColl, juiceTran.items[:])
// 		} else if juiceGiftTran.length > 0 {
// 			err = wfj.BatchInsert(ctx, wfj.JuiceGiftColl, juiceGiftTran.items[:])
// 		} else if juicePurchasedTran.length > 0 {
// 			err = wfj.BatchInsert(ctx, wfj.JuicePurchasedColl, juicePurchasedTran.items[:])
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }

func (wfj *WaterFallJob) DealAndInsert(ctx context.Context, data any) error {
	var err error
	freeCountTran := new(TranContainer)
	juiceTran := new(TranContainer)
	juiceGiftTran := new(TranContainer)
	juicePurchasedTran := new(TranContainer)
	if data, ok := data.([]*models.BankAccountTransaction); ok {
		for i := range data {
			switch data[i].BankType {
			case models.FreeCount:
				err = wfj.CollClassify(ctx, wfj.FreeCountColl, data[i], freeCountTran)
			case models.Juice:
				err = wfj.CollClassify(ctx, wfj.JuiceColl, data[i], juiceTran)
			case models.JuiceGift:
				err = wfj.CollClassify(ctx, wfj.JuiceGiftColl, data[i], juiceGiftTran)
			case models.JuicePurchased:
				err = wfj.CollClassify(ctx, wfj.JuicePurchasedColl, data[i], juicePurchasedTran)
			}
			if err != nil {
				return err
			}
		}
	}
	if freeCountTran.length > 0 {
		err = wfj.BatchInsert(ctx, wfj.FreeCountColl, freeCountTran.items[:])
	} else if juiceTran.length > 0 {
		err = wfj.BatchInsert(ctx, wfj.JuiceColl, juiceTran.items[:])
	} else if juiceGiftTran.length > 0 {
		err = wfj.BatchInsert(ctx, wfj.JuiceGiftColl, juiceGiftTran.items[:])
	} else if juicePurchasedTran.length > 0 {
		err = wfj.BatchInsert(ctx, wfj.JuicePurchasedColl, juicePurchasedTran.items[:])
	}
	if err != nil {
		return err
	}
	return nil
}

func (wfj *WaterFallJob) Read(ctx context.Context) (any, error) {
	col := wfj.WfDb.Collection(wfj.LogColl)
	logs := []*models.Log{}
	filter := models.M{
		"opTypes": models.M{
			"$nin": []string{"initiateAccount"},
		},
	}
	// cur, err := col.Find(ctx, filter, options.Find().SetLimit(400).SetSkip(400))
	cur, err := col.Find(ctx, filter)
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
