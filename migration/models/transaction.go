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
	IsMockUser          bool               `bson:"isMockUser"`
	IsDeleted           bool               `bson:"isDeleted"`
	CreatedAt           time.Time          `bson:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt"`
}

type BankAccountTransactionSource struct {
	TransactionID string `bson:"transactionId"`
	Amount        int64  `bson:"amount"`
}

type CombinationChildren struct {
	BankType      BankType `bson:"bankType"`
	TransactionID string   `bson:"transactionId"`
	Amount        int64    `bson:"amount"`
	RestAmount    int64    `bson:"restAmount"`
}

type BankAccountTransaction struct {
	*Transaction         `bson:",inline"`
	BankType             BankType                                `bson:"bankType"` // 银行帐户类型,juice,points等
	Reason               string                                  `bson:"reason"`
	Amount               int64                                   `bson:"amount"`
	RestAmount           int64                                   `bson:"restAmount"`           // 可使用的数量, 默认restAmount=amount, 对于reload,capture,refund,transferIn等可用的钱, 需要记录已经使用的数量
	SourceTransactionIDs map[string]BankAccountTransactionSource `bson:"sourceTransactionIds"` // 需要可追述, 记录该次交易的上级
	ChildTransactionIDs  map[string]BankAccountTransactionSource `bson:"childTransactionIds"`  // 记录钱的去处
	Operation            BankAccountOperation                    `bson:"operation"`
	TransID              string                                  `bson:"transID"`
}
