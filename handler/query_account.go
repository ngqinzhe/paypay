package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ngqinzhe/paypay/dal/dao"
	"github.com/ngqinzhe/paypay/errs"
)

type QueryAccountHandler struct {
	accountDao *dao.AccountDao
}

func NewQueryAccountHandler(accountDao *dao.AccountDao) *QueryAccountHandler {
	return &QueryAccountHandler{
		accountDao: accountDao,
	}
}

func (h *QueryAccountHandler) Handle(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIdReq := c.Param("account_id")
		if accountIdReq == "" {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: "no account id provided",
			})
			return
		}
		accountId, err := strconv.ParseInt(accountIdReq, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpErrorResponse{
				Message: errs.ErrInvalidRequest.Error(),
				Error:   errors.New("invalid accountId").Error(),
			})
			return
		}
		account, err := h.accountDao.QueryAccount(ctx, accountId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpErrorResponse{
				Message: errs.ErrServerErr.Error(),
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, queryAccountResponse{
			Account: *account,
			Message: "query account success",
		})
	}
}
