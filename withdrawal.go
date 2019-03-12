package main

import "go.mongodb.org/mongo-driver/bson/primitive"

var withdrawalCollection = DB.Collection("withdrawals")

type Withdrawal struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Amount       int64                `bson:"amount"`
	Transactions []DepositTransaction `bson:"transactions"`
	Status       bool                 `bson:"status"`
}

type WithdrawalTransaction struct {
	ID         primitive.ObjectID `bson:"_id"`
	Amount     int64              `bson:"amount"`
	Comment    string             `bson:"comment"`
	ContractID string             `bson:"contractID"`
	Status     bool               `bson:"status"`
}