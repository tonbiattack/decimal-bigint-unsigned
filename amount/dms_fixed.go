package amount

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// ImportedRowFixed は修正後の DMS インポート行。
// DMS のマッピング設定を文字列に変更することで float64 変換を回避する。
type ImportedRowFixed struct {
	Amount string
}

// SafeMapFromDMSFixed は文字列で受け取った金額を decimal に変換する。
// float64 を経由しないため精度が失われない。
func SafeMapFromDMSFixed(row ImportedRowFixed) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(row.Amount)
	if err != nil {
		return decimal.Zero, fmt.Errorf("decimal への変換に失敗しました: %w", err)
	}
	if d.IsNegative() {
		return decimal.Zero, errors.New("amount は正の値である必要があります")
	}
	return d, nil
}
