package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(tx Transaction) error {
	if err := r.db.Create(&tx).Error; err != nil {
		return fmt.Errorf("insert に失敗しました: %w", err)
	}
	return nil
}

func (r *Repository) FindAll() ([]Transaction, error) {
	var txs []Transaction
	if err := r.db.Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("取得に失敗しました: %w", err)
	}
	return txs, nil
}
