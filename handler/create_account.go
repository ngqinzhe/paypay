package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngqinzhe/paypay/dal/dao"
	"github.com/ngqinzhe/paypay/dal/model"
	"github.com/ngqinzhe/paypay/errs"
)

type CreateAccountHandler struct {
	accountDao *dao.AccountDao
}

func NewCreateAccountHandler(accountDao *dao.AccountDao) *CreateAccountHandler {
	return &CreateAccountHandler{
		accountDao: accountDao,
	}
}

func (h *CreateAccountHandler) Handle(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &createAccountRequest{}
		if err := c.BindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: errs.ErrInvalidRequest.Error(),
				Error:   err.Error(),
			})
			return
		}
		if req.AccountId == 0 || req.InitialBalance == "" {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: "invalid account_id or initial_balance",
			})
			return
		}
		account := &model.Account{
			AccountId: req.AccountId,
			Balance:   req.InitialBalance,
		}
		if err := h.accountDao.CreateAccount(ctx, account); err != nil {
			c.JSON(http.StatusInternalServerError, httpErrorResponse{
				Message: errs.ErrServerErr.Error(),
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, createAccountResponse{
			Account: *account,
			Message: "account created successfully",
		})
	}
}
