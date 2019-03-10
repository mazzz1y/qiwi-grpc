package main

import "C"
import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pb "qiwi/protobuf"
)

func (s *Server) CreateOrUpdateAccount(ctx context.Context, in *pb.CreateOrUpdateAccountRequest) (*pb.CreateOrUpdateAccountResponse, error) {
	a := Account{Token: in.Token}
	a, err := a.CreateOrUpdate()
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateAccountResponse{
		ContractID:             a.ContractID,
		OperationLimit:         a.OperationLimit,
		MaxAllowableBalance:    a.MaxAllowableBalance,
		OperationLimitPerMonth: a.OperationLimitPerMonth,
		Balance:                a.Balance,
		Blocked:                a.Blocked},
		nil
}

func (s *Server) ListAccounts(ctx context.Context, in *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{
		ContractIDs: Account{}.ListString()},
		nil
}

func (s *Server) GetAccountBalances(ctx context.Context, in *pb.GetAccountBalancesRequest) (*pb.GetAccountBalancesResponse, error) {
	b, err := Account{ContractID: in.ContractID}.GetBalance()
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountBalancesResponse{
		Balance: b},
		nil
}

func (s *Server) Deposit(ctx context.Context, in *pb.DepositRequest) (*pb.DepositResponse, error) {
	deposit, err := Deposit{Amount: in.Amount}.Create()

	var amounts []int64
	var links []string
	var comments []string
	var contractIDs []string

	for _, d := range deposit.Transactions {
		amounts = append(amounts, d.Amount)
		links = append(links, d.Link)
		comments = append(comments, d.Comment)
		contractIDs = append(contractIDs, d.ContractID)
	}
	if err != nil {
		return nil, err
	}
	return &pb.DepositResponse{
		Id:          deposit.ID.Hex(),
		ContractIDs: contractIDs,
		Amounts:     amounts,
		Links:       links,
		Comments:    comments},
		nil
}

func (s *Server) DepositCheck(ctx context.Context, in *pb.DepositCheckRequest) (*pb.DepositCheckResponse, error) {
	objID, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, err //todo readable
	}
	deposit, err := Deposit{ID: objID}.Check()
	if err != nil {
		return nil, err
	}

	var amounts []int64
	var links []string
	var comments []string
	var contractIDs []string
	var statuses []bool

	for _, d := range deposit.Transactions {
		amounts = append(amounts, d.Amount)
		links = append(links, d.Link)
		comments = append(comments, d.Comment)
		contractIDs = append(contractIDs, d.ContractID)
		statuses = append(statuses, d.Status)
	}
	if err != nil {
		return nil, err
	}
	return &pb.DepositCheckResponse{
		Id:          deposit.ID.Hex(),
		Status:      deposit.Status,
		ContractIDs: contractIDs,
		Amounts:     amounts,
		Links:       links,
		Comments:    comments,
		Statuses:    statuses},
		nil
}
