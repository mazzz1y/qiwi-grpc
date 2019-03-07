package main

import "C"
import (
	"context"
	pb "qiwi/protobuf"
)

func (s *Server) AddAccount(ctx context.Context, in *pb.AddAccountRequest) (*pb.AddAccountResponse, error) {
	a := Account{Token: in.Token, ContractID: in.ContractID}
	a, err := a.Create()
	if err != nil {
		return nil, err
	}
	return &pb.AddAccountResponse{ContractID: a.ContractID, OperationLimit: a.OperationLimit}, nil
}

func (s *Server) ListAccounts(ctx context.Context, in *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{ContractIDs: Account{}.ListString()}, nil
}

func (s *Server) GetAccountBalances(ctx context.Context, in *pb.GetAccountBalancesRequest) (*pb.GetAccountBalancesResponse, error) {
	b, err := Account{ContractID: in.ContractID}.GetBalance()
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountBalancesResponse{Balance: b}, nil
}

func (s *Server) SendMoneyToQiwi(ctx context.Context, in *pb.SendMoneyToQiwiRequest) (*pb.SendMoneyToQiwiResponse, error) {
	statusCode, err := Account{ContractID: in.SenderContractID}.SendMoneyToQiwi(in.Amount, in.ReceiverContractID)
	if err != nil {
		return nil, err
	}
	return &pb.SendMoneyToQiwiResponse{Status: statusCode}, nil
}

func (s *Server) GetPaymentLinks(ctx context.Context, in *pb.GetPaymentLinksRequest) (*pb.GetPaymentLinksResponse, error) {
	pl := GetPaymentLinks(in.Amount, "")
	return &pb.GetPaymentLinksResponse{PaymentLinks: pl}, nil
}
