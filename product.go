package jewerly

import "errors"

// Products
type CreateProductInput struct {
	Titles        MultiLanguageInput `json:"titles" binding:"required"`
	Descriptions  MultiLanguageInput `json:"descriptions" binding:"required"`
	CurrentPrice  float32            `json:"current_price" binding:"required"`
	PreviousPrice float32            `json:"previous_price"`
	Code          string             `json:"code" binding:"required"`
	ImageIds      []int              `json:"image_ids" binding:"required"`
	CategoryId    Category           `json:"category_id" binding:"required"`
}

func (i CreateProductInput) Validate() error {
	return i.CategoryId.Validate()
}

type MultiLanguageInput struct {
	English   string `json:"english" binding:"required"`
	Russian   string `json:"russian" binding:"required"`
	Ukrainian string `json:"ukrainian" binding:"required"`
}

// Categories

type Category int

func (c Category) Validate() error {
	_, ok := Categories[c]
	if !ok {
		return errors.New("invalid category")
	}

	return nil
}

const (
	CategoryRings = iota + 1
	CategoryBracelets
	CategoryPendants
	CategoryEarring
	CategoryNecklaces
)

var (
	Categories = map[Category]bool{
		CategoryRings:     true,
		CategoryBracelets: true,
		CategoryPendants:  true,
		CategoryEarring:   true,
		CategoryNecklaces: true,
	}
)
