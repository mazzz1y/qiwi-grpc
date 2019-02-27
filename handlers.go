package main

import "C"
import (
	"context"

	pb "qiwi/protobuf"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddAccount(ctx context.Context, in *pb.AddAccountRequest) (*pb.AddAccountResponse, error) {
	a := Account{Token: in.Token}
	a, err := a.Create()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			err.Error())
	}
	return &pb.AddAccountResponse{ContractID: a.ContractID}, nil
}

func (s *Server) ListAccounts(ctx context.Context, in *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{ContractIDs: Account{}.List()}, nil
}

func (s *Server) GetAccountBalances(ctx context.Context, in *pb.GetAccountBalancesRequest) (*pb.GetAccountBalancesResponse, error) {
	b, err := Account{ContractID: in.ContractID}.Balance()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			err.Error())
	}
	return &pb.GetAccountBalancesResponse{Balances: b}, nil
}

func (s *Server) SendMoneyToQiwi(ctx context.Context, in *pb.SendMoneyToQiwiRequest) (*pb.SendMoneyToQiwiResponse, error) {
	statusCode, err := Account{ContractID: in.SenderContractID}.SendMoneyToQiwi(in.Amount, in.ReceiverContractID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			err.Error())
	}

	return &pb.SendMoneyToQiwiResponse{Status: statusCode}, nil
}