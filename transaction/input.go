package transaction

import "funding-app/users"

type GetTransactionCampaignInput struct {
	ID   int `uri:"id" binding:"required"`
	User users.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       users.User
}
