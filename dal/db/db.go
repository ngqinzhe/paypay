package db

import (
	"context"

	"github.com/ngqinzhe/paypay/consts"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient interface {
	GetCollection(collectionName string) *mongo.Collection
	Close(ctx context.Context)
	StartSession(ctx context.Context) (mongo.Session, error)
}

const (
	driverStr  = "mongodb+srv://ngqinzhe2:qBL9Q0JbqYyu9nhA@cluster0.wrjlydx.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	sessionKey = "paypay_session"
)

func Init() MongoDBClient {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(driverStr))
	if err != nil {
		panic(err)
	}
	return &mongoDB{
		client: client,
	}
}

func WithSession(ctx context.Context, collectionName string) *mongo.Collection {
	v := ctx.Value(sessionKey)
	session, ok := v.(mongo.Session)
	if !ok {
		return nil
	}
	return session.Client().Database(consts.MongoDbName).Collection(collectionName)
}

type mongoDB struct {
	client *mongo.Client
}

func (c *mongoDB) GetCollection(collectionName string) *mongo.Collection {
	return c.client.Database(consts.MongoDbName).Collection(collectionName)
}

func (c *mongoDB) StartSession(ctx context.Context) (mongo.Session, error) {
	session, err := c.client.StartSession()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, sessionKey, session)
	return session, nil
}

func (c *mongoDB) Close(ctx context.Context) {
	if err := c.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
