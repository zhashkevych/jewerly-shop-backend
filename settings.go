package jewerly

import "gopkg.in/guregu/null.v3"

type UpdateSettingsInput struct {
	HeaderText         *MultiLanguageInput `json:"header_text"`
	TextBlock1         *MultiLanguageInput `json:"text_block_1"`
	TextBlock2         *MultiLanguageInput `json:"text_block_2"`
	CustomerService    *MultiLanguageInput `json:"customer_service"`
	AboutUs            *MultiLanguageInput `json:"about_us"`
	PrivacyPolicy      *MultiLanguageInput `json:"privacy_policy"`
	TermsAndConditions *MultiLanguageInput `json:"terms_and_conditions"`
	ShippingAndReturns *MultiLanguageInput `json:"shipping_and_returns"`
	SocialLinks        []SocialLink        `json:"social_links"`
	FrontPageImages    []FrontPageImage    `json:"front_page_images"`
	MinimalOrderPrice  null.Int            `json:"minimal_order_price"`
	StoreName          null.String         `json:"store_name"`
	LogoImageId        null.Int            `json:"logo_image_id"`
}

type SocialLink struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	ImageId string `json:"image_id"`
}

type FrontPageImage struct {
	Id      int `json:"id"`
	ImageId int `json:"image_id"`
}
