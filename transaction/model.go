package transaction

import "time"

// Transaction は transactions テーブルのモデル。
// AmountUnsafe は BIGINT UNSIGNED にマッピングされており、型設計が不正な例。
// AmountSafe は DECIMAL(35,20) にマッピングされており、正しい例。
type Transaction struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	AmountUnsafe  uint64    `gorm:"column:amount_unsafe"`
	AmountSafe    string    `gorm:"column:amount_safe"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}
