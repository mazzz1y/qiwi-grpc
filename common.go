package main

const nedeedPercentForPayment = 30

//todo check if wallet is not blocked for payment
func GetPaymentLinks(amount int64, comment string) ([]string) {
	accounts := Account{}.List()
	var readyAccounts []Account
	var links []string
	comment = "321" //todo
	for _, acc := range accounts {
		//BUG need to check limits instead of balance
		if acc.Balance/100*nedeedPercentForPayment > amount {
			readyAccounts = append(readyAccounts, acc)
		}
	}

	for _, acc := range readyAccounts {
		links = append(links, acc.GetPaymentLink(int(amount), comment))
	}
	return links
}
