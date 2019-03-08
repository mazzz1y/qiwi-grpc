package main

const nedeedPercentForPayment = 30
const maxAmountPerWallet = 15000


//todo check max wallet balance
//todo
//todo check if wallet is not blocked for payment
func GetPaymentLinks(amount int64, comment string) ([]string) {
	accounts := Account{}.List()
	var links []string
	comment = "321" //todo
	for _, acc := range accounts {
		walAvalPay := acc.MaxAllowableBalance - acc.Balance
		if amount == 0 {
			break
		}
		if walAvalPay > amount {
			links = append(links, acc.GetPaymentLink(int(amount), comment))
			amount-=amount
			continue
		} else if (walAvalPay < amount) && (walAvalPay/100*nedeedPercentForPayment < amount) {
			links = append(links, acc.GetPaymentLink(int(walAvalPay), comment))
			amount-=walAvalPay
			continue
		}
	}

	return links
}
