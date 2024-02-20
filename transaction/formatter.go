package transaction

import "time"

type TransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatTransactionCampaign(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}
func FormatTransactionCampaigns(transactions []Transaction) []TransactionFormatter {
	if len(transactions) == 0 {
		return []TransactionFormatter{}
	}
	var transactionFormatter []TransactionFormatter
	for _, transaction := range transactions {
		formatter := FormatTransactionCampaign(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}
	return transactionFormatter
}
