package amount

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// UnsafeConvertToUint64 は decimal を float64 を経由して uint64 に変換する。
// 精度欠損・アンダーフローが発生する危険な実装例。
func UnsafeConvertToUint64(d decimal.Decimal) uint64 {
	f, _ := d.Float64()
	// 負数を uint64 にキャストすると最大値に化ける
	return uint64(int64(f))
}

// SafeConvertToString は decimal を文字列として返す。
// 精度を失わずに扱うための安全な実装例。
func SafeConvertToString(d decimal.Decimal) string {
	return d.String()
}

// SafeConvertToDecimal は文字列から decimal に変換する。
func SafeConvertToDecimal(s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, fmt.Errorf("decimal への変換に失敗しました: %w", err)
	}
	return d, nil
}
