package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryFixed struct {
	db *gorm.DB
}

func NewRepositoryFixed(db *gorm.DB) *RepositoryFixed {
	return &RepositoryFixed{db: db}
}

func (r *RepositoryFixed) Insert(tx TransactionFixed) error {
	if err := r.db.Create(&tx).Error; err != nil {
		return fmt.Errorf("insert に失敗しました: %w", err)
	}
	return nil
}

func (r *RepositoryFixed) FindAll() ([]TransactionFixed, error) {
	var txs []TransactionFixed
	if err := r.db.Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("取得に失敗しました: %w", err)
	}
	return txs, nil
}
