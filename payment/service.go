package payment

import (
	"funding-app/users"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"strconv"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user users.User) (string, error)
}

func NewService() *service {
	return &service{}
}
func (s *service) GetPaymentURL(transaction Transaction, user users.User) (string, error) {
	// 1. Initiate Snap client
	var sn = snap.Client{}
	sn.New("SB-Mid-server-e7Nl5WC8JnmX04ioqKqW3PDf", midtrans.Sandbox)

	// 2. Initiate Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Name,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, err := sn.CreateTransaction(snapReq)
	if err != nil {
		return "", err
	}
	return snapResp.RedirectURL, nil
}
