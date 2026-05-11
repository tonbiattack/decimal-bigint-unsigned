package transaction

import (
	"time"

	"github.com/shopspring/decimal"
)

// TransactionFixed は修正後の transactions テーブルのモデル。
// 金額を shopspring/decimal で扱い、DECIMAL(35,20) カラムにマッピングする。
type TransactionFixed struct {
	ID           uint64          `gorm:"primaryKey;autoIncrement"`
	AmountFixed  decimal.Decimal `gorm:"column:amount_fixed;type:decimal(35,20)"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime"`
}

func (TransactionFixed) TableName() string {
	return "transactions"
}
