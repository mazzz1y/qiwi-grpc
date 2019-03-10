package main

func getBalanceSum() int64 {
	var b int64
	accounts := Account{}.List()
	for _, acc := range accounts {
		b += acc.Balance
	}
	return b
}
