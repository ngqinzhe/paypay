package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngqinzhe/paypay/dal/dao"
	"github.com/ngqinzhe/paypay/dal/db"
	"github.com/ngqinzhe/paypay/dal/model"
	"github.com/ngqinzhe/paypay/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTransactionHandler struct {
	dbClient       db.MongoDBClient
	transactionDao *dao.TransactionDao
	accountDao     *dao.AccountDao
}

func NewCreateTransactionHandler(dbClient db.MongoDBClient, transactionDao *dao.TransactionDao, accountDao *dao.AccountDao) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		dbClient:       dbClient,
		transactionDao: transactionDao,
		accountDao:     accountDao,
	}
}

func (h *CreateTransactionHandler) Handle(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &model.Transaction{}
		if err := c.BindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: errs.ErrInvalidRequest.Error(),
				Error:   err.Error(),
			})
			return
		}
		if !isValidCreateTransactionRequest(req) {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: errs.ErrInvalidRequest.Error(),
				Error:   "transaction request parameters invalid",
			})
		}

		// update time
		req.Time = time.Now().Unix()
		srcBal, destBal, err := h.createTransaction(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpErrorResponse{
				Message: errs.ErrServerErr.Error(),
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, &transactionResponse{
			SourceAccountId:      req.SourceAccountId,
			SourceAccountBalance: srcBal,
			DestinationAccountId: req.DestinationAccountId,
			DestinationBalance:   destBal,
			Message:              fmt.Sprintf("transaction of $%s success", req.Amount),
		})
	}
}

func isValidCreateTransactionRequest(req *model.Transaction) bool {
	if req.Amount == "" {
		return false
	}
	if req.SourceAccountId == 0 {
		return false
	}
	if req.DestinationAccountId == 0 {
		return false
	}
	if req.SourceAccountId == req.DestinationAccountId {
		return false
	}
	return true
}

func (h *CreateTransactionHandler) createTransaction(ctx context.Context, transaction *model.Transaction) (string, string, error) {
	var (
		err                error
		srcUpdatedBalance  string
		destUpdatedBalance string
	)
	transferAmt, err := strconv.ParseFloat(transaction.Amount, 64)
	if err != nil {
		fmt.Println("error here")
		return "", "", err
	}
	// query source account whether there are enough funds
	sourceAccount, err := h.accountDao.QueryAccount(ctx, transaction.SourceAccountId)
	if err != nil {
		return "", "", err
	}

	srcBalance, _ := strconv.ParseFloat(sourceAccount.Balance, 64)
	if !isBalanceEnoughForTransfer(srcBalance, transferAmt) {
		return "", "", errors.New("transaction amount exceeded balance")
	}
	// query destination account
	destinationAccount, err := h.accountDao.QueryAccount(ctx, transaction.DestinationAccountId)
	if err != nil {
		return "", "", err
	}

	destBalance, _ := strconv.ParseFloat(destinationAccount.Balance, 64)
	// perform the transaction and update balance
	srcUpdatedBalance = updateBalance(srcBalance, transferAmt, true)
	destUpdatedBalance = updateBalance(destBalance, transferAmt, false)

	session, err := h.dbClient.StartSession(ctx)
	if err != nil {
		return "", "", err
	}
	// abort transaction if err
	defer func(ctx context.Context, err error) (string, string, error) {
		if err != nil {
			session.AbortTransaction(ctx)
			return "", "", err
		}
		return srcUpdatedBalance, destUpdatedBalance, nil
	}(ctx, err)

	if err := session.StartTransaction(); err != nil {
		return "", "", err
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := h.accountDao.UpdateAccount(ctx, sourceAccount.AccountId, bson.D{{"balance", srcUpdatedBalance}}); err != nil {
			return err
		}
		if err := h.accountDao.UpdateAccount(ctx, destinationAccount.AccountId, bson.D{{"balance", destUpdatedBalance}}); err != nil {
			return err
		}
		if err := h.transactionDao.CreateTransactionRecord(ctx, transaction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", "", err
	}
	return "", "", session.CommitTransaction(ctx)
}

func isBalanceEnoughForTransfer(balance, transferAmount float64) bool {
	return balance >= transferAmount
}

// updateBalance will return updated balance, if isSrc is true, the transfer amount will be deducted, vice versa
func updateBalance(balance, transferAmount float64, isSrc bool) string {
	var updated float64
	if isSrc {
		updated = balance - transferAmount
	} else {
		updated = balance + transferAmount
	}
	return strconv.FormatFloat(updated, 'f', -1, 64)
}
