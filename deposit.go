package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strconv"
)

const maxAmountPerWallet = 15000

var depositCollection = DB.Collection("deposit")

type DepositTransaction struct {
	Amount  int64  `bson:"amount"`
	Comment string `bson:"comment"`
	Link    string `bson:"link"`
}

type Deposit struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Amount       int64                `bson:"amount"`
	Transactions []DepositTransaction `bson:"transactions"`
	Status       string               `bson:"status"`
}

func (d Deposit) Create() (Deposit, error) {
	tr, err := getPaymentLinks(d.Amount)
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

//todo check if wallet is not blocked for payment
func getPaymentLinks(amount int64) ([]DepositTransaction, error) {
	accounts := Account{}.List()
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
				Amount:  amount,
				Comment: comment,
				Link:    acc.GetPaymentLink(int(amount), comment)})
			amount -= amount
		} else if availSumForDeposit <= amount {
			depositTransactions = append(depositTransactions, DepositTransaction{
				Amount:  availSumForDeposit,
				Comment: comment,
				Link:    acc.GetPaymentLink(int(availSumForDeposit), comment)})
			amount -= availSumForDeposit
		}
	}
	return depositTransactions, nil
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
