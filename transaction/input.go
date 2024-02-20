package transaction

import "funding-app/users"

type GetTransactionCampaignInput struct {
	ID   int `uri:"id" binding:"required"`
	User users.User
}
