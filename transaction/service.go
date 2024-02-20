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

	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
