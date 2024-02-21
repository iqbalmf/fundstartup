package payment

import (
	"funding-app/campaign"
	"funding-app/transaction"
	"funding-app/users"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"strconv"
)

type service struct {
	transactionRepo transaction.Repository
	campaignRepo    campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user users.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(transRepository transaction.Repository, campaignRepository campaign.Repository) *service {
	return &service{transRepository, campaignRepository}
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

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error {
	transactionId, _ := strconv.Atoi(input.OrderID)
	transactionTemp, err := s.transactionRepo.GetTransactionByID(transactionId)
	if err != nil {
		return err
	}
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transactionTemp.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transactionTemp.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" ||
		input.TransactionStatus == "cancel" {
		transactionTemp.Status = "cancelled"
	} else {
		transactionTemp.Status = "cancelled"
	}
	updatedTransaction, err := s.transactionRepo.UpdateTransaction(transactionTemp)
	if err != nil {
		return err
	}
	campaignTemp, err := s.campaignRepo.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status == "paid" {
		campaignTemp.BackerCount = campaignTemp.BackerCount + 1
		campaignTemp.CurrentAmount = campaignTemp.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepo.UpdateCampaign(campaignTemp)
		if err != nil {
			return err
		}
	}
	return nil
}
