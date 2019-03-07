package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var trCollection = DB.Collection("transactions")

type Transaction struct {
	Amount int64 `bson:"amount"`
	Type string `bson:"type"`
	ContractID string `bson:"contractID"`
}

func (t Transaction) Create() (Transaction, error) {
	_, err := trCollection.InsertOne(context.TODO(), t)
	if err != nil {
		return t, status.Errorf(codes.Internal, err.Error())
	}
	return t, nil
}

func (t Transaction) Close() error {
	r, err := trCollection.DeleteOne(context.TODO(), t)
	if err != nil || r.DeletedCount < 1 {
		return status.Errorf(codes.Internal, "transaction not found")
	}
	return nil
}

