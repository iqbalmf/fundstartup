package transaction

import (
	"errors"
	"funding-app/campaign"
	"funding-app/payment"
	"strconv"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}
type Service interface {
	GetTransactionByCampaignID(input GetTransactionCampaignInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	SaveTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
	GetAllTransaction() ([]Transaction, error)
}

func NewService(repository Repository, campaignRepo campaign.Repository, payment payment.Service) *service {
	return &service{repository: repository, campaignRepository: campaignRepo, paymentService: payment}
}

func (s *service) GetTransactionByCampaignID(input GetTransactionCampaignInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}

	transaction, err := s.repository.GetTransactionByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transaction, err := s.repository.GetTransactionByUserID(userID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (s *service) SaveTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = ""

	newTransaction, err := s.repository.SaveTransaction(transaction)
	if err != nil {
		return newTransaction, err
	}
	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.UpdateTransaction(newTransaction)
	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transactionId, _ := strconv.Atoi(input.OrderID)
	transactionTemp, err := s.repository.GetTransactionByID(transactionId)
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
	updatedTransaction, err := s.repository.UpdateTransaction(transactionTemp)
	if err != nil {
		return err
	}
	campaignTemp, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status == "paid" {
		campaignTemp.BackerCount = campaignTemp.BackerCount + 1
		campaignTemp.CurrentAmount = campaignTemp.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.UpdateCampaign(campaignTemp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetAllTransaction() ([]Transaction, error) {
	transactions, err := s.repository.FindAllTransaction()
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
