package jewerly

type TextBlock struct {
	ID int `json:"id"`
	MultiLanguageInput
}

type Settings struct {
	Images     []Image     `json:"images"`
	TextBlocks []TextBlock `json:"text-blocks"`
}
