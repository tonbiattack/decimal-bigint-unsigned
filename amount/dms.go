package amount

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// ImportedRow は DMS 経由でインポートされたデータの構造体。
// DMS が連携元の decimal カラムを float64 に変換するため、フィールドが float64 になっている。
type ImportedRow struct {
	Amount float64
}

// DestinationRow はインポート先システムの構造体。
// 金額を uint64 で持つ設計になっている。
type DestinationRow struct {
	Amount uint64
}

// UnsafeMapFromDMS は DMS インポート行をインポート先の struct に変換する。
// float64 を uint64 に直接キャストするため、値が範囲外のとき結果は実装依存になる。
func UnsafeMapFromDMS(row ImportedRow) DestinationRow {
	return DestinationRow{
		// float64 → uint64 の直接キャスト。
		// 値が uint64 の範囲外（負値など）のとき、結果は実装依存になる。
		Amount: uint64(row.Amount),
	}
}

// SafeMapFromDMS は DMS インポート行を精度を保ったまま decimal に変換する。
func SafeMapFromDMS(row ImportedRow) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(fmt.Sprintf("%.20f", row.Amount))
	if err != nil {
		return decimal.Zero, fmt.Errorf("decimal への変換に失敗しました: %w", err)
	}
	return d, nil
}
