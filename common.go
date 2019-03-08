package main

const nedeedPercentForPayment = 30
const maxAmountPerWallet = 15000


//todo check max wallet balance
//todo implement maxAmountPerWallet
//todo check if wallet is not blocked for payment
func GetPaymentLinks(amount int64, comment string) []string {
	accounts := Account{}.List()
	var links []string
	comment = "321" //todo
	for _, acc := range accounts {
		if amount == 0 {
			break
		}
		walAvalPay := acc.MaxAllowableBalance - acc.Balance
		if walAvalPay > maxAmountPerWallet {
			walAvalPay = maxAmountPerWallet
		}
		if walAvalPay >= amount {
			links = append(links, acc.GetPaymentLink(int(amount), comment))
			amount-=amount
		} else if walAvalPay <= amount {
			links = append(links, acc.GetPaymentLink(int(walAvalPay), comment))
			amount -= walAvalPay
		}
	}

	return links
}
