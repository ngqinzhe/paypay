package dao

import (
	"context"

	"github.com/ngqinzhe/paypay/consts"
	"github.com/ngqinzhe/paypay/dal/db"
	"github.com/ngqinzhe/paypay/dal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionDao struct {
	collection *mongo.Collection
}

func NewTransactionDao(dbClient db.MongoDBClient) *TransactionDao {
	return &TransactionDao{
		collection: dbClient.GetCollection(consts.MongoDbCollectionName_Transactions),
	}
}

func (t *TransactionDao) CreateTransactionRecord(ctx context.Context, transaction *model.Transaction) error {
	collection := t.collection

	sessionCollection := db.WithSession(ctx, consts.MongoDbCollectionName_Transactions)
	if sessionCollection != nil {
		collection = sessionCollection
	}

	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
