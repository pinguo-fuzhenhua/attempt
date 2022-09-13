package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TS             time.Time          `bson:"ts"`                 //时间戳
	AppName        string             `bson:"appName"`            //产品名称
	DataChannel    string             `bson:"dataChannel"`        //数据渠道
	UserID         string             `bson:"userID,omitempty"`   //用户id
	DeviceID       string             `bson:"deviceID,omitempty"` //设备id
	OpType         string             `bson:"opType"`             //操作类型
	OriginalAmount int64              `bson:"originalAmount"`     //操作前数量
	Amount         int64              `bson:"amount"`             //操作数量
	Comment        string             `bson:"comment,omitempty"`  //描述
	TransID        string             `bson:"transID"`            //保证幂等的ID
}

type PointsOperateType string

const (
	// 购买
	POTPurchase PointsOperateType = "purchase"
	// 赠送
	POTGift PointsOperateType = "gift"
	// 消耗
	POTConsume PointsOperateType = "consume"
	// 系统清理
	POTSysClear PointsOperateType = "sysclear"
	// 转账
	POTTransfer PointsOperateType = "transfer"
)

func (l *Log) ToTransaction() *BankAccountTransaction {
	bankType, opType, err := l.getBankTypeAndOpType()
	if err != nil {
		return nil
	}

	return &BankAccountTransaction{
		Transaction: &Transaction{
			ID:                  l.ID,
			OriginTransactionID: l.buildTransID(PointsOperateType(l.getOperation()), l.ID.Hex(), "", l.UserID),
			TransactionID:       "",
			UserID:              l.UserID,
			DeviceID:            l.DeviceID,
			Comment:             l.Comment,
			Scope:               strings.ToLower(l.AppName),
			CreatedAt:           l.TS,
			UpdatedAt:           l.TS,
		},
		BankType:  BankType(bankType),
		Amount:    l.Amount,
		Reason:    l.OpType,
		Operation: BankAccountOperation(opType),
		TransID:   l.TransID,
	}
}

func (l *Log) getBankTypeAndOpType() (string, string, error) {
	var (
		bankType string
		opType   string
	)
	switch strings.ToLower(l.OpType) {
	case "helpconvert", "fps", "speed", "rsmb", "clear":
		switch l.OpType {
		case "helpconvert":
			opType = string(Reload)
		default:
			opType = string(Sale)
		}
		bankType = l.DataChannel
	case "purchasepreset":
		bankType = string(Juice)
		opType = string(TransferOut)
	case "purchasejuice":
		bankType = string(JuicePurchased)
		opType = string(Reload)
	case "gift", "newcomer":
		bankType = string(JuiceGift)
		opType = string(Reload)
	default:
		return "", "", errors.New("无效数据")
	}
	return bankType, opType, nil
}

func (l *Log) getOperation() string {
	operation := ""
	optype := strings.ToLower(l.OpType)
	if strings.Contains(optype, "clear") {
		operation = string(POTSysClear)
	} else if strings.Contains(optype, "preset") {
		operation = string(POTTransfer)
	} else if strings.Contains(optype, "gift") {
		operation = string(POTGift)
	} else if strings.Contains(optype, "purchase") || strings.Contains(optype, "convert") {
		operation = string(POTPurchase)
	} else {
		operation = string(POTConsume)
	}
	return operation
}

func (l *Log) buildTransID(operateType PointsOperateType, id, sid string, userID string) string {
	if sid != "" {
		id = id + ":" + sid
	}
	return string(operateType) + "-" + id + "-" + userID
}
