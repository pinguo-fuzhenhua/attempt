package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Scope       string             `bson:"scope"`
	MockUserIDs []string           `bson:"mockUserIDs"`
	DeviceID    string             `bson:"deviceId"`
	IsDeleted   bool               `bson:"isDeleted"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

type BankResult struct {
	BankType         BankType            `bson:"bankType"`
	Amount           int64               `bson:"amount"`
	AuthorizedAmount int64               `bson:"authorizedAmount"`
	Offset           *primitive.ObjectID `bson:"offset"` // 快照时的最后一条事务ID
}

type UserSnapshot struct {
	ID        ID                       `bson:"_id"`
	Scope     string                   `bson:"scope"`
	IsMock    bool                     `bson:"isMock"`
	UserID    string                   `bson:"userId"`
	Banks     map[BankType]*BankResult `bson:"banks"`
	UpdatedAt time.Time                `bson:"updatedAt"`
	CreatedAt time.Time                `bson:"createdAt"`
}

type MiddleUser struct {
	ID         ID        `json:"_id"`
	SnapshotID ID        `json:"snapshotID"`
	Account    []string  `json:"account"`
	Amount     []int64   `json:"amount"`
	UpdateAt   time.Time `json:"updateAt"`
	CreatedAt  time.Time `json:"createAt"`
	Scope      string    `json:"scope"`
}

func (mu *MiddleUser) ToSnapshot() *UserSnapshot {
	banks := make(map[BankType]*BankResult)
	for i := range mu.Account {
		banks[BankType(mu.Account[i])] = &BankResult{
			BankType:         BankType(mu.Account[i]),
			Amount:           mu.Amount[i],
			AuthorizedAmount: 0,
			Offset:           &primitive.NilObjectID,
		}
	}
	return &UserSnapshot{
		ID:        mu.SnapshotID,
		Scope:     strings.ToLower(mu.Scope),
		IsMock:    false,
		UserID:    mu.ID.Hex(),
		Banks:     banks,
		CreatedAt: mu.CreatedAt,
	}
}
