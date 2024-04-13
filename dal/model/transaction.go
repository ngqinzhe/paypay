package model

type Transaction struct {
	SourceAccountId      int64  `json:"source_account_id" bson:"source_account_id"`
	DestinationAccountId int64  `json:"destination_account_id" bson:"destination_account_id"`
	Amount               string `json:"amount" bson:"amount"`
	Time                 int64  `json:"time,omitempty" bson:"time"`
}
