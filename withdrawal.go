package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var withdrawalCollection = DB.Collection("withdrawals")

type Withdrawal struct {
	ID                 primitive.ObjectID      `bson:"_id"`
	Amount             int64                   `bson:"amount"`
	Transactions       []WithdrawalTransaction `bson:"transactions"`
	ReceiverContractID string                  `bson:"receiverContractID"`
	Status             string                  `bson:"status"`
	Comment            string                  `bson:"comment"`
}

type WithdrawalTransaction struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Amount             int64              `bson:"amount"`
	Comment            string             `bson:"comment"`
	SenderContractID   string             `bson:"senderContractID"`
	ReceiverContractID string             `bson:"receiverContractID"`
	Status             string             `bson:"status"`
}

func (w Withdrawal) Create() (Withdrawal, error) {
	trs, err := w.generateTransactions()
	if err != nil {
		return w, err
	}

	w.ID = primitive.NewObjectID()
	w.Transactions = trs
	w.Status = "failed"

	//todo to-lower
	for _, tr := range trs {
		log.Println(tr.Status)
		if tr.Status != "Accepted" {
			w.Status = tr.Status
		} else {
			w.Status = "accepted"
		}
	}

	_, err = withdrawalCollection.InsertOne(context.TODO(), w)
	if err != nil {
		return w, status.Errorf(codes.Internal, err.Error())
	}
	return w, nil
}

func (w Withdrawal) generateTransactions() ([]WithdrawalTransaction, error) {
	accounts := Account{}.List()
	withdrawSum := getBalanceSum()
	defaultStatus := "pending"
	amount := w.Amount
	comment := w.Comment
	var withdrawalTransactions []WithdrawalTransaction

	if amount > withdrawSum {
		return withdrawalTransactions, status.Errorf(codes.OutOfRange, "not enough sum for withdraw")
	}
	for _, acc := range accounts {
		if amount == 0 {
			break
		}
		availSumForWithdraw := acc.Balance
		if availSumForWithdraw >= amount {
			withdrawalTransactions = append(withdrawalTransactions, WithdrawalTransaction{
				ID:                 primitive.NewObjectID(),
				Amount:             amount,
				Comment:            comment,
				SenderContractID:   acc.ContractID,
				ReceiverContractID: w.ReceiverContractID,
				Status:             defaultStatus})
			amount -= amount
		} else if availSumForWithdraw <= amount {
			withdrawalTransactions = append(withdrawalTransactions, WithdrawalTransaction{
				ID:                 primitive.NewObjectID(),
				Amount:             availSumForWithdraw,
				Comment:            comment,
				SenderContractID:   acc.ContractID,
				ReceiverContractID: w.ReceiverContractID,
				Status:             defaultStatus})
			amount -= availSumForWithdraw
		}
	}

	for i, tr := range withdrawalTransactions {
		trCode, err := tr.makeTransaction()
		withdrawalTransactions[i].Status = trCode
		if err != nil {
			return withdrawalTransactions, err
		}
	}
	return withdrawalTransactions, nil
}

//todo rename
func (w WithdrawalTransaction) makeTransaction() (string, error) {
	acc := Account{ContractID: w.SenderContractID}

	trCode, err := acc.SendMoneyToQiwi(w.Amount, w.ReceiverContractID, w.Comment)

	if err != nil {
		return trCode, status.Errorf(codes.Aborted, "withdraw error")
	}

	w.Status = trCode

	return trCode, nil
}

func getBalanceSum() int64 {
	var b int64
	accounts := Account{}.List()
	for _, acc := range accounts {
		b += acc.Balance
	}
	return b
}
