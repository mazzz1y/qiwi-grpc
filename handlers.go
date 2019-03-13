package main

import "C"
import (
	"context"
	pb "qiwi/protobuf"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *Server) DepositCreate(ctx context.Context, in *pb.DepositCreateRequest) (*pb.DepositCreateResponse, error) {
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
	return &pb.DepositCreateResponse{
		Id:          deposit.ID.Hex(),
		ContractIDs: contractIDs,
		Amounts:     amounts,
		Links:       links,
		Comments:    comments},
		nil
}

func (s *Server) DepositClose(ctx context.Context, in *pb.DepositCloseRequest) (*pb.DepositCloseResponse, error) {
	objID, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, err //todo readable
	}
	deposit, err := Deposit{ID: objID}.Close()
	if err != nil {
		return nil, err
	}
	return &pb.DepositCloseResponse{
		Id:     deposit.ID.Hex(),
		Status: deposit.Status,
	},
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

func (s *Server) WithdrawalCreate(ctx context.Context, in *pb.WithdrawalCreateRequest) (*pb.WithdrawalCreateResponse, error) {
	withdrawal, err := Withdrawal{Amount: in.Amount, ReceiverContractID: in.ContractID, Comment: in.Comment}.Create()

	var amounts []int64
	var comments []string
	var receiverContractIDs []string
	var senderContractIDs []string
	var statuses []string

	for _, tr := range withdrawal.Transactions {
		amounts = append(amounts, tr.Amount)
		comments = append(comments, tr.Comment)
		receiverContractIDs = append(receiverContractIDs, tr.ReceiverContractID)
		senderContractIDs = append(senderContractIDs, tr.SenderContractID)
		statuses = append(statuses, tr.Status)
	}
	if err != nil {
		return nil, err
	}
	return &pb.WithdrawalCreateResponse{
		Id:                  withdrawal.ID.Hex(),
		Status:              withdrawal.Status,
		ReceiverContractIDs: receiverContractIDs,
		SenderContractIDs:   senderContractIDs,
		Amounts:             amounts,
		Statuses:            statuses,
		Comments:            comments},
		nil
}
