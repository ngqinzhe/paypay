package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ngqinzhe/paypay/dal/dao"
	"github.com/ngqinzhe/paypay/dal/db"
	"github.com/ngqinzhe/paypay/handler"
)

func main() {
	ctx := context.Background()
	s := newServer(gin.Default(), db.Init())
	s.Init(ctx)
	defer s.Shutdown(ctx)
}

type server struct {
	router *gin.Engine
	db     db.MongoDBClient
}

func newServer(router *gin.Engine, db db.MongoDBClient) *server {
	return &server{
		router: router,
		db:     db,
	}
}

func (s *server) Init(ctx context.Context) {
	s.router.Use(cors.Default())
	s.initRoutes(ctx)
}

func (s *server) Shutdown(ctx context.Context) {
	s.db.Close(ctx)
}

func (s *server) initRoutes(ctx context.Context) {
	// get
	// s.router.GET("/accounts", handler.NewRootHandler().Handle(ctx))
	accountDao := dao.NewAccountDao(s.db)
	transactionDao := dao.NewTransactionDao(s.db)
	// account creation
	s.router.POST("/accounts", handler.NewCreateAccountHandler(accountDao).Handle(ctx))
	// account query
	s.router.GET("/accounts/:account_id", handler.NewQueryAccountHandler(accountDao).Handle(ctx))
	// create transaction
	s.router.POST("/transactions", handler.NewCreateTransactionHandler(s.db, transactionDao, accountDao).Handle(ctx))

	s.router.Run("localhost:3000")
}
