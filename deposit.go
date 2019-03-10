package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strconv"
)

const maxAmountPerWallet = 15000

var depositCollection = DB.Collection("deposits")

type Deposit struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Amount       int64                `bson:"amount"`
	Transactions []DepositTransaction `bson:"transactions"`
	Status       bool                 `bson:"status"`
}

type DepositTransaction struct {
	ID         primitive.ObjectID `bson:"_id"`
	Amount     int64              `bson:"amount"`
	Comment    string             `bson:"comment"`
	ContractID string             `bson:"contractID"`
	Link       string             `bson:"link"`
	Status     bool               `bson:"status"`
}

func (d Deposit) Create() (Deposit, error) {
	tr, err := d.generateTransactions()
	if err != nil {
		return d, err
	}

	d.Transactions = tr
	d.ID = primitive.NewObjectID()

	_, err = depositCollection.InsertOne(context.TODO(), d)
	if err != nil {
		return d, status.Errorf(codes.Internal, err.Error())
	}
	return d, nil
}

func (d Deposit) Check() (Deposit, error) {
	err := depositCollection.FindOne(context.TODO(), bson.M{"_id": d.ID}).Decode(&d)
	if len(d.Transactions) == 0 {
		return d, status.Errorf(codes.NotFound, "transactions not found in db")
	}
	d, err = d.updateTransactionsStatus()
	if err != nil {
		return d, err
	}
	return d, nil
}

//todo check if wallet is not blocked for payment
//todo block available sum
func (d Deposit) generateTransactions() ([]DepositTransaction, error) {
	accounts := Account{}.List()
	amount := d.Amount
	var depositTransactions []DepositTransaction

	if getAvailableDepositSum() < amount {
		return depositTransactions, status.Errorf(codes.OutOfRange, "not enough deposit sum")
	}
	for _, acc := range accounts {
		comment := strconv.FormatInt(rand.Int63(), 10)
		if amount == 0 {
			break
		}
		availSumForDeposit := acc.MaxAllowableBalance - acc.Balance
		if availSumForDeposit > maxAmountPerWallet {
			availSumForDeposit = maxAmountPerWallet
		}
		if availSumForDeposit >= amount {
			depositTransactions = append(depositTransactions, DepositTransaction{
				ID: primitive.NewObjectID(),
				Amount:     amount,
				Comment:    comment,
				ContractID: acc.ContractID,
				Link:       acc.GetPaymentLink(int(amount), comment),
				Status:     false})
			amount -= amount
		} else if availSumForDeposit <= amount {
			depositTransactions = append(depositTransactions, DepositTransaction{
				ID: primitive.NewObjectID(),
				Amount:     availSumForDeposit,
				Comment:    comment,
				ContractID: acc.ContractID,
				Link:       acc.GetPaymentLink(int(availSumForDeposit), comment),
				Status:     false})
			amount -= availSumForDeposit
		}
	}
	return depositTransactions, nil
}

func (d Deposit) updateTransactionsStatus() (Deposit, error) {
	const days = 1
	for i, tr := range d.Transactions {
		trStatus, err := Account{ContractID: tr.ContractID}.CheckTransactionIsSuccess(days, tr.Comment, tr.Amount)
		if err != nil {
			return d, err
		}
		d.Transactions[i].Status = trStatus
		// if transaction is false, return deposit without change
		if !trStatus {
			return d, nil
		}
	}
	d.Status = true
	_, err := depositCollection.ReplaceOne(context.TODO(), bson.M{"_id": d.ID}, &d)
	if err != nil {
		return d, status.Errorf(codes.Internal, "failed to update transaction status")
	}
	return d, nil
}

func getAvailableDepositSum() int64 {
	var sum int64
	accounts := Account{}.List()
	for _, acc := range accounts {
		max := acc.MaxAllowableBalance - acc.Balance
		if max > maxAmountPerWallet {
			sum += maxAmountPerWallet
		} else {
			sum += max
		}
	}
	return sum
}
