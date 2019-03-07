package main

import "log"

const nedeedPercentForPayment = 30

//todo check if wallet is not blocked for payment
func GetPaymentLinks(amount int64, comment string) ([]string) {
	accounts := Account{}.List()
	var links []string
	comment = "321" //todo
	var total int64
	for _, acc := range accounts {
		sum := acc.MaxAllowableBalance - acc.Balance
		log.Println(acc.MaxAllowableBalance)
		log.Println(acc.Balance)
		log.Println(sum)
		if  amount > sum/100*nedeedPercentForPayment   {
			if amount - (total + sum) <= 0 {
				links = append(links, acc.GetPaymentLink(int(amount-total), ""))
				log.Println(links)
				break
			}
			links = append(links, acc.GetPaymentLink(int(amount), ""))
		}
	}

	return links
}
