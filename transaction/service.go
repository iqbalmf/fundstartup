package transaction

import (
	"errors"
	"funding-app/campaign"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}
type Service interface {
	GetTransactionByCampaignID(input GetTransactionCampaignInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	SaveTransaction(input CreateTransactionInput) (Transaction, error)
}

func NewService(repository Repository, campaignRepo campaign.Repository) *service {
	return &service{repository: repository, campaignRepository: campaignRepo}
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
	return newTransaction, nil
}
