package transaction

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Transaction interface {
	Do(ctx context.Context, fn ...func(ctx context.Context) error) error
}

type trans struct {
	client *mongo.Client

	opt *options.TransactionOptions
}

func NewTransaction(c *mongo.Database) Transaction {
	return &trans{
		client: c.Client(),
		opt:    options.Transaction().SetReadPreference(readpref.Primary()),
	}
}

func (t *trans) Do(ctx context.Context, fn ...func(ctx context.Context) error) error {
	var rfn func(context.Context) error
	switch len(fn) {
	case 0:
		return nil
	case 1:
		rfn = fn[0]
	default:
		rfn = func(ctx context.Context) error {
			for i := range fn {
				if err := fn[i](ctx); err != nil {
					return err
				}
			}
			return nil
		}
	}

	return t.client.UseSession(ctx, func(sc mongo.SessionContext) error {
		_, err := sc.WithTransaction(sc, func(sessCtx mongo.SessionContext) (interface{}, error) {
			err := rfn(sessCtx)
			return nil, err
		}, t.opt)
		return err
	})
}
