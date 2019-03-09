package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/url"
	"qiwi/client"
	"strconv"
	"time"
)

const (
	paymentMethod = 99  // per-account payment method code
	currency      = 643 // RUB currency
)

var (
	accCollection       = DB.Collection("accounts")
	maxAllowableBalance = map[string]int64{
		"SIMPLE":   15000,
		"VERIFIED": 60000,
		"FULL":     600000,
	}
	operationLimitPerMonth = map[string]int64{
		"SIMPLE":   40000,
		"VERIFIED": 200000,
		"FULL":     999999999, // hardcoded unlimited
	}
	operationLimit = map[string]int64{
		"SIMPLE":   15000,
		"VERIFIED": 60000,
		"FULL":     500000,
	}
)

type Account struct {
	Token                  string `bson:"token"`
	ContractID             string `bson:"contractID"`
	OperationLimit         int64  `bson:"operationLimit"`
	MaxAllowableBalance    int64  `bson:"maxAllowableBalance"`
	OperationLimitPerMonth int64  `bson:"operationLimitPerMonth"`
	Balance                int64  `bson:"balance"`
	Blocked                bool   `bson:"blocked"`
}

// Create or update account with a new token
func (a Account) Create() (Account, error) {
	a, err := a.init()
	if status.Code(err) != 0 {
		return a, err
	}
	// update account if already exist
	count, err := accCollection.CountDocuments(context.TODO(), bson.M{"contractID": a.ContractID})
	if count > 0 {
		_, err = accCollection.ReplaceOne(context.TODO(), bson.M{"contractID": a.ContractID}, &a)
		if err != nil {
			return a, status.Errorf(codes.Internal, err.Error())
		}
		return a, nil
	}

	_, err = accCollection.InsertOne(context.TODO(), a)
	if err != nil {
		return a, status.Errorf(codes.Internal, err.Error())
	}

	err = SetUniqueIndex(accCollection, "contractID")
	if err != nil {
		return a, status.Errorf(codes.Internal, err.Error())
	}

	return a, nil
}

// Return stored accounts
func (Account) List() (accounts []Account) {
	cur, err := accCollection.Find(context.TODO(), bson.M{})
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
		accounts = append(accounts, a)
	}
	return accounts
}

// Return contractID for stored accounts
func (Account) ListString() (accountsString []string) {
	accounts := Account{}.List()
	for _, a := range accounts {
		accountsString = append(accountsString, a.ContractID)
	}
	return accountsString
}

// Return balance with each currencies for account
func (a Account) GetBalance() (int64, error) {
	c, err := a.client()
	if status.Code(err) != 0 {
		return 0, err
	}
	balances, err := c.Balance.Current()
	if status.Code(err) != 0 {
		return 0, status.Errorf(codes.Internal, err.Error())
	}
	for _, b := range balances.Accounts {
		if b.BankAlias == "QIWI" && b.Currency == currency {
			return int64(b.Balance.Amount), nil
		}
	}
	return 0, status.Errorf(codes.Unavailable, "balance unavailable")
}

// Send money to a Qiwi Wallet, only RUB
func (a Account) SendMoneyToQiwi(amount int64, receiverContractID string) (string, error) {
	c, err := a.client()
	if status.Code(err) != 0 {
		return "", err
	}
	payment, err := c.Cards.Payment(paymentMethod, float64(amount), receiverContractID)

	if payment.Transaction.State.Code != "Accepted" {
		err := errors.New("transaction declined")
		return "", status.Errorf(codes.Aborted, err.Error())
	}
	return payment.Transaction.State.Code, nil
}

// Get user-friendly payment link
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
	baseURL := paymentFormURL + strconv.Itoa(paymentMethod) + "?" + v.Encode()
	return baseURL
}

// Check that account has enough limits for making out payment
func (a Account) IsReadyForMakePayment(amount int64) (bool, error) {
	const days = 31

	if amount > a.OperationLimit {
		return false, nil
	}
	var paymentSum float64
	paymentHist, err := a.getPaymentsHistory(days)
	if err != nil {
		return false, err
	}
	for _, payment := range paymentHist {
		if payment.Sum.Currency == currency {
			paymentSum += payment.Sum.Amount
		}
	}
	if int64(paymentSum) > a.OperationLimitPerMonth {
		return false, nil
	}
	return true, nil
}
//todo implement ability to refresh wallet
func (a Account) init() (Account, error) {
	c, _ := a.client()
	profile, err := c.Profile.Current()
	if err != nil {
		return a, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if profile.ContractInfo.Blocked {
		return a, status.Errorf(codes.PermissionDenied, err.Error())
	}
	log.Println(a)
	balance, err := a.GetBalance()
	log.Println(err)
	if err != nil {
		return a, status.Errorf(codes.Internal, "failed to get balance from API")
	}

	identLevel := getQiwiIdentificationLevel(profile)

	a.ContractID = strconv.FormatInt(profile.AuthInfo.PersonID, 10)
	a.OperationLimit = operationLimit[identLevel]
	a.MaxAllowableBalance = maxAllowableBalance[identLevel]
	a.OperationLimitPerMonth = operationLimitPerMonth[identLevel]
	a.OperationLimit = operationLimit[identLevel]
	a.Balance = balance

	return a, nil
}

func (a Account) getPaymentsHistory(days int) ([]client.Txn, error) {
	const rows = 50
	c, err := a.client()
	if status.Code(err) != 0 {
		return nil, err
	}
	v := url.Values{}
	v.Set("startDate", time.Now().AddDate(0, 0, -days).Format(time.RFC3339))
	v.Set("endDate", time.Now().Format(time.RFC3339))
	pr, err := c.Payments.History(rows, v)
	if err != nil {
		return pr.Data, status.Errorf(codes.Internal, "getting history error")
	}
	return pr.Data, nil
}

// Return qiwi-client
func (a Account) client() (*client.Client, error) {
	_ = accCollection.FindOne(context.TODO(), bson.M{"contractID": a.ContractID}).Decode(&a)
	c := client.New(a.Token)
	c.SetWallet(a.ContractID)
	if Config.Debug {
		c.Debug = true
	}
	return c, nil
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
