package dao

import (
	"context"

	"github.com/ngqinzhe/paypay/consts"
	"github.com/ngqinzhe/paypay/dal/db"
	"github.com/ngqinzhe/paypay/dal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountDao struct {
	collection *mongo.Collection
}

func NewAccountDao(dbClient db.MongoDBClient) *AccountDao {
	return &AccountDao{
		collection: dbClient.GetCollection(consts.MongoDbCollectionName_Accounts),
	}
}

func (a *AccountDao) CreateAccount(ctx context.Context, account *model.Account) error {
	collection := a.collection

	sessionCollection := db.WithSession(ctx, consts.MongoDbCollectionName_Accounts)
	if sessionCollection != nil {
		collection = a.collection
	}
	if _, err := collection.InsertOne(ctx, account); err != nil {
		return err
	}
	return nil
}

func (a *AccountDao) UpdateAccount(ctx context.Context, accountId int64, updates bson.D) error {
	collection := a.collection

	sessionCollection := db.WithSession(ctx, consts.MongoDbCollectionName_Accounts)
	if sessionCollection != nil {
		collection = a.collection
	}
	_, err := collection.UpdateOne(ctx,
		bson.D{{"account_id", accountId}},
		bson.D{{"$set", updates}})
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDao) QueryAccount(ctx context.Context, accountId int64) (*model.Account, error) {
	singleResult := a.collection.FindOne(ctx, bson.D{{"account_id", accountId}})
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}
	account := &model.Account{}
	if err := singleResult.Decode(account); err != nil {
		return nil, err
	}
	return account, nil
}
