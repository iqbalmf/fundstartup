package campaign

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"'`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	var imageUrl = ""
	if len(campaign.CampaignImages) > 0 {
		imageUrl = campaign.CampaignImages[0].FileName
	}
	formatter := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		ImageUrl:         imageUrl,
	}
	return formatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	if len(campaigns) == 0 {
		return []CampaignFormatter{}
	}
	var campaignsFormatter []CampaignFormatter
	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}
	return campaignsFormatter
}
