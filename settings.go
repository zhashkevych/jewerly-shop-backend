package jewerly

import (
	"errors"
	"gopkg.in/guregu/null.v3"
)

type TextBlock struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	MultiLanguageInput
}

type UpdateTextBlockInput struct {
	Name null.String         `json:"name"`
	Text *MultiLanguageInput `json:"text"`
}

func (i UpdateTextBlockInput) Validate() error {
	if !i.Name.Valid && i.Text == nil {
		return errors.New("no update values")
	}
	return nil
}

type HomepageImage struct {
	ID  int    `json:"id" db:"id"`
	URL string `json:"url" db:"url"`
}
type Settings struct {
	Images     []HomepageImage `json:"images"`
	TextBlocks []TextBlock     `json:"text-blocks"`
}
