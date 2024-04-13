package model

// data model for Account
type Account struct {
	AccountId int64  `json:"account_id" bson:"account_id"`
	Balance   string `json:"balance" bson:"balance"`
}
