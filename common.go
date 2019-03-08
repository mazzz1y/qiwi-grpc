package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxAmountPerWallet = 15000

//todo check max wallet balance
//todo check if wallet is not blocked for payment
func GetPaymentLinks(amount int64, comment string) ([]string, error) {
	accounts := Account{}.List()
	var links []string
	if getAvailableDepositSum() < amount {
		return links, status.Errorf(codes.OutOfRange, "not enough deposit sum")
	}
	comment = "321" //todo
	for _, acc := range accounts {
		if amount == 0 {
			break
		}
		availSumForDeposit := acc.MaxAllowableBalance - acc.Balance
		if availSumForDeposit > maxAmountPerWallet {
			availSumForDeposit = maxAmountPerWallet
		}
		if availSumForDeposit >= amount {
			links = append(links, acc.GetPaymentLink(int(amount), comment))
			amount-=amount
		} else if availSumForDeposit <= amount {
			links = append(links, acc.GetPaymentLink(int(availSumForDeposit), comment))
			amount -= availSumForDeposit
		}
	}
	return links, nil
}

func getAvailableDepositSum() int64 {
	var b int64
	accounts := Account{}.List()
	for _, acc := range accounts {
		b += acc.MaxAllowableBalance - acc.Balance
	}
	return b
}

func getBalanceSum() int64 {
	var b int64
	accounts := Account{}.List()
	for _, acc := range accounts {
		b += acc.Balance
	}
	return b
}
