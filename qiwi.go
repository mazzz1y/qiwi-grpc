package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/url"
	"qiwi/client"
	"strconv"
)

// per-account payment method code
const paymentMethod = 99

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
	Token      string             `bson:"token"`
	ContractID string              `bson:"contractID"`
	Type	   string			  `bson:"type"` //FULL,VERIFIED,SIMPLE
}

// Create or update account with a new token
func (a Account) Create() (Account, error) {
	c, _ := a.client()
	profile, err := c.Profile.Current()
	if err != nil {
		return a, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if profile.ContractInfo.Blocked {
		return a, status.Errorf(codes.PermissionDenied, err.Error())
	}

	a.ContractID = strconv.FormatInt(profile.AuthInfo.PersonID, 10)
	a.Type = getQiwiIdentificationLevel(profile)

	// update account if already exist
	count, err := collection.CountDocuments(context.TODO(), bson.M{"contractID": a.ContractID})
	if count > 0 {
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"contractID": a.ContractID}, &a)
		if err != nil {
			return a, status.Errorf(codes.Internal, err.Error())
		}
		return a, nil
	}

	_, err = collection.InsertOne(context.TODO(), a)
	if err != nil {
		return a, status.Errorf(codes.Internal, err.Error())
	}

	err = setUniqueIndex("contractID", collection)
	if err != nil {
		return a, status.Errorf(codes.Internal, err.Error())
	}

	return a, nil
}

// Return contractID for stored accounts
func (Account) List() (list []string) {
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var a Account
		err := cur.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, a.ContractID)
	}
	return list
}

// Return balance with each currencies for account
func (a Account) Balance() (map[string]float64, error) {
	m := map[string]float64{}
	c, err := a.client()
	if status.Code(err) != 0 {
		return m, err
	}
	b, err := c.Balance.Current()
	if status.Code(err) != 0 {
		return m, status.Errorf(codes.Internal, err.Error())
	}
	for _, cur := range b.Accounts {
		m[currencies[cur.Currency]] = cur.Balance.Amount
	}
	return m, nil
}

// Send money to a Qiwi Wallet, only RUB
func (a Account) SendMoneyToQiwi(amount float64, receiverContractID string) (string, error) {
	c, err := a.client()
	if status.Code(err) != 0 {
		return "", err
	}
	payment, err := c.Cards.Payment(paymentMethod, amount, receiverContractID)

	if payment.Transaction.State.Code != "Accepted" {
		err := errors.New("transaction declined")
		return "", status.Errorf(codes.Aborted, err.Error())
	}
	return payment.Transaction.State.Code, nil
}

func (a Account) GetPaymentLink(amount int, comment string) string {
	const amountFraction = 0
	const paymentFormURL = "https://qiwi.com/payment/form/"
	const currency = 643

	v := url.Values{}
	v.Set("extra['account']", a.ContractID)
	v.Set("currency", strconv.Itoa(currency))
	v.Set("amountInteger", strconv.Itoa(amount))
	v.Set("amountFraction", strconv.Itoa(amountFraction))
	v.Set("extra['comment']", comment)
	// block input
	v.Set("blocked[0]", "account")
	v.Set("blocked[1]", "comment")
	baseURL := paymentFormURL + strconv.Itoa(paymentMethod) + "?" +v.Encode()
	return baseURL
}

// Return qiwi-client if account exist in db
func (a Account) client() (*client.Client, error) {
	err := collection.FindOne(context.TODO(), bson.M{"contractID": a.ContractID}).Decode(&a)
	c := client.New(a.Token)
	if Config.Debug {
		c.Debug = true
	}
	if err != nil {
		return c, status.Errorf(codes.NotFound, "account not found in db")
	}
	return c, nil
}

// Create unique index if number of rows = 1
func setUniqueIndex(rowName string, collection *mongo.Collection) error {
	allDocCount, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	if allDocCount == 1 {
		_, err = collection.Indexes().CreateOne(
			context.TODO(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{rowName, bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			return err
		}
	}
	return nil
}

// Get identification level for qiwi bank
func getQiwiIdentificationLevel(p client.ProfileResponse) string {
	for _, bank := range p.ContractInfo.IdentificationInfo {
		if bank.BankAlias == "QIWI" {
			return bank.IdentificationLevel
		}
	}
	return ""
}
