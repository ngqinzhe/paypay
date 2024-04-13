package handler

import "github.com/ngqinzhe/paypay/dal/model"

type createAccountRequest struct {
	AccountId      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type createAccountResponse struct {
	model.Account
	Message string `json:"message"`
}

type queryAccountResponse struct {
	model.Account
	Message string `json:"message"`
}

type transactionResponse struct {
	SourceAccountId      int64  `json:"source_account_id"`
	SourceAccountBalance string `json:"source_account_balance"`
	DestinationAccountId int64  `json:"destination_account_id"`
	DestinationBalance   string `json:"destination_account_balance"`
	Message              string `json:"message"`
}

type httpErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
