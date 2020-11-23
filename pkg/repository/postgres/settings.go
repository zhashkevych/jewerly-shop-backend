package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type SettingsRepository struct {
	db *sqlx.DB
}

func NewSettingsRepository(db *sqlx.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetImages() ([]jewerly.HomepageImage, error) {
	var images []jewerly.HomepageImage
	err := r.db.Select(&images, fmt.Sprintf("SELECT hi.id, i.url FROM %s hi INNER JOIN %s i ON hi.image_id = i.id",
		homepageImagesTable, imagesTable))
	return images, err
}

func (r *SettingsRepository) CreateImage(imageID int) error {
	query := fmt.Sprintf("INSERT INTO %s (image_id) VALUES ($1)", homepageImagesTable)
	_, err := r.db.Exec(query, imageID)
	return err
}

func (r *SettingsRepository) UpdateImage(id, imageID int) error {
	query := fmt.Sprintf("UPDATE %s SET image_id = $1 WHERE id = $2", homepageImagesTable)
	_, err := r.db.Exec(query, imageID, id)
	return err
}

func (r *SettingsRepository) GetTextBlocks() ([]jewerly.TextBlock, error) {
	var textBlocks []jewerly.TextBlock
	err := r.db.Select(&textBlocks, fmt.Sprintf("SELECT tb.id, tb.name, ml.english, ml.russian, ml.ukrainian FROM %s tb INNER JOIN %s ml ON tb.text_id=ml.id",
		textBlocksTable, multiLanguageTextTable))
	return textBlocks, err
}

func (r *SettingsRepository) GetTextBlockById(id int) (jewerly.TextBlock, error) {
	var textBlock jewerly.TextBlock
	err := r.db.Get(&textBlock, fmt.Sprintf("SELECT tb.id, tb.name, ml.english, ml.russian, ml.ukrainian FROM %s tb INNER JOIN %s ml ON tb.text_id=ml.id WHERE tb.id=$1",
		textBlocksTable, multiLanguageTextTable), id)
	return textBlock, err
}

func (r *SettingsRepository) CreateTextBlock(block jewerly.TextBlock) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (english, russian, ukrainian) VALUES ($1, $2, $3) RETURNING id",
		multiLanguageTextTable)
	row := tx.QueryRow(query, block.English, block.Russian, block.Ukrainian)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (name, text_id) VALUES ($1, $2)", textBlocksTable)
	_, err = tx.Exec(query, block.Name, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SettingsRepository) UpdateTextBlock(id int, block jewerly.UpdateTextBlockInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if block.Text != nil {
		query := fmt.Sprintf("UPDATE %s ml SET english=$1, russian=$2, ukrainian=$3 FROM %s tb WHERE tb.text_id = ml.id AND tb.id=$4",
			multiLanguageTextTable, textBlocksTable)
		_, err := tx.Exec(query, block.Text.English, block.Text.Russian, block.Text.Ukrainian, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if block.Name.Valid {
		query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", textBlocksTable)
		_, err := tx.Exec(query, block.Name, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
