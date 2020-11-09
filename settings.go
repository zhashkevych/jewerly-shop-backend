package jewerly

import "gopkg.in/guregu/null.v3"

const (
	PageCustomerService = "customer-service"
	PageAboutUs         = "about-us"
)

type SetHomepageInput struct {
	ProductIDs []int               `json:"product_ids"`
	ImageId    null.Int            `json:"image_id"`
	TextBlock1 *MultiLanguageInput `json:"text_block_1"`
	TextBlock2 *MultiLanguageInput `json:"text_block_2"`
}

type HomepageSettings struct {
	Products   ProductsList       `json:"products"`
	ImageURL   string             `json:"image_url"`
	TextBlock1 MultiLanguageInput `json:"text_block_1"`
	TextBlock2 MultiLanguageInput `json:"text_block_2"`
}

type TextPage struct {
	Text string `json:"text"`
}
