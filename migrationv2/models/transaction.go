package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	BankType             string
	BankAccountOperation string
)

const (
	Sale        BankAccountOperation = "sale"
	Reload      BankAccountOperation = "reload"
	TransferIn  BankAccountOperation = "transferIn"
	TransferOut BankAccountOperation = "transferOut"

	Juice          BankType = "juice"
	JuiceGift      BankType = "juice_gift"
	JuicePurchased BankType = "juice_purchased"
	FreeCount      BankType = "free_count"
)

type Transaction struct {
	ID                  primitive.ObjectID `bson:"_id"`
	TransactionID       string             `bson:"transactionId"`
	OriginTransactionID string             `bson:"originTransactionId"`
	Scope               string             `bson:"scope"`
	UserID              string             `bson:"userId"`
	Comment             string             `bson:"comment"`
	DeviceID            string             `bson:"deviceId"`
	Type                string             `bson:"type"`
	CreatedAt           time.Time          `bson:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt"`
}

type BankAccountTransaction struct {
	*Transaction `bson:",inline"`
	BankType     BankType             `bson:"bankType"` // 银行帐户类型,juice,points等
	Reason       string               `bson:"reason"`
	Amount       int64                `bson:"amount"`
	Operation    BankAccountOperation `bson:"operation"`
	TransID      string               `bson:"transID"`
}
