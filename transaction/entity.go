package transaction

import (
	"funding-app/campaign"
	"funding-app/users"
	"github.com/leekchan/accounting"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	User       users.User
	Campaign   campaign.Campaign
	PaymentURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t Transaction) FormattedTime() string {
	dateTimeLayout := "2006-01-02 15:04:05"
	formattedTime := t.UpdatedAt.Format(dateTimeLayout)
	return formattedTime
}

func (c Transaction) AmontFormatIDR() string {
	acc := accounting.Accounting{
		Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ",",
	}
	return acc.FormatMoney(c.Amount)
}
