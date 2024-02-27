package campaign

import (
	"funding-app/users"
	"github.com/leekchan/accounting"
	"time"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             users.User
}

func (c Campaign) GoalAmountFormatIDR() string {
	acc := accounting.Accounting{
		Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ",",
	}
	return acc.FormatMoney(c.GoalAmount)
}

func (c Campaign) CurrentAmontFormatIDR() string {
	acc := accounting.Accounting{
		Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ",",
	}
	return acc.FormatMoney(c.CurrentAmount)
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
