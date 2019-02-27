package main

import (
	"context"
	"errors"
	"github.com/zhuharev/qiwi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

var (
	collection = DB.Collection("accounts")
	currencies = map[int]string{
		978: "EUR",
		398: "KZT",
		756: "CHF",
		972: "TJS",
		840: "USD",
		980: "UAH",
		643: "RUB",
	}

)

type Account struct {
	Token string             `bson:"token"`
	ContractID int64              `bson:"contractID"`
}

// Create or update account with a new token
func (a Account) Create() (Account, error) {
	c, _ := a.Client()
	profile, err := c.Profile.Current()
	if err != nil {
		return a, err
	}
	if profile.ContractInfo.Blocked != false {
		return a, errors.New("Account is blocked")
	}

	a.ContractID = profile.AuthInfo.PersonID

	// update account if already exist
	count, err := collection.CountDocuments(context.TODO(), bson.M{"contractID": a.ContractID})
	if count > 0 {
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"contractID": a.ContractID}, &a)
		log.Println(a)
		return a, err
	}

	_, err = collection.InsertOne(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}

	// create unique index if number of accounts = 1
	allDocCount, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if allDocCount == 1 {
		_, err = collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"contractID", bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			log.Fatal(err)
		}
	}
	return a, nil
}

// Return contractID for stored accounts
func (Account) List() (list []int64) {
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var a Account
		err := cur.Decode(&a)
		if err != nil { log.Fatal(err) }
		list = append(list, a.ContractID)
	}
	return list
}

// Return balance with each currencies for account
func (a Account) Balance() (map[string]float64, error) {
	m := map[string]float64{}
	c, err := a.Client()
	if err != nil {
		return m, err
	}
	b, err := c.Balance.Current()
	if err != nil {
		return m, err
	}
	for _, cur := range b.Accounts {
		m[currencies[cur.Currency]] = cur.Balance.Amount
	}
		return m, err
	}

// Return qiwi-client if account exist in db
func (a Account) Client() (*qiwi.Client, error) {
	err := collection.FindOne(context.TODO(), bson.M{"contractID": a.ContractID}).Decode(&a)
	client := qiwi.New(a.Token)
	return client, err
}

// Send money to a Qiwi Wallet, only RUB
func (a Account) SendMoneyToQiwi(amount float64, receiverContractID string) (string, error) {
	// per-account payment method code
	const paymentMethod = 99
	c, err := a.Client()
	if err != nil {
		return "", err
	}
	d, err := c.Cards.Payment(paymentMethod, amount, receiverContractID)
	return d.Transaction.State.Code, err
}
