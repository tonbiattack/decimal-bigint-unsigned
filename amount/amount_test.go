package amount_test

import (
	"testing"

	"github.com/tonbiattack/decimal-bigint-unsigned/amount"
)

// uint64 の最大値
const maxUint64 = uint64(18446744073709551615)

// TestUnsafeMap_小数部が切り捨てられる は、DMS が float64 で渡した小数値が uint64 変換で切り捨てられることを示す。
func TestUnsafeMap_小数部が切り捨てられる(t *testing.T) {
	// DB の decimal カラムが DMS 経由で float64 になった値
	row := amount.ImportedRow{Amount: 12345678.123456789}
	got := amount.UnsafeMap(row)
	// 小数部が切り捨てられ、整数部だけ残る
	if got.Amount != 12345678 {
		t.Errorf("got %d, want 12345678", got.Amount)
	}
}

// TestUnsafeMap_極小値がゼロになる は、0.000...001 のような値がゼロに化けることを示す。
func TestUnsafeMap_極小値がゼロになる(t *testing.T) {
	row := amount.ImportedRow{Amount: 0.000000000000000001}
	got := amount.UnsafeMap(row)
	if got.Amount != 0 {
		t.Errorf("got %d, want 0", got.Amount)
	}
}

// TestUnsafeMap_負数がuint64最大値に化ける は、DMS が負数を float64 で渡したとき uint64 への変換で最大値になることを示す。
// Go の仕様では符号付き整数から符号なし整数への変換は 2^N を法とした値になるため、
// -1 (int64) を uint64 に変換すると 18446744073709551615 になる。
func TestUnsafeMap_負数がuint64最大値に化ける(t *testing.T) {
	row := amount.ImportedRow{Amount: -1}
	got := amount.UnsafeMap(row)
	if got.Amount != maxUint64 {
		t.Errorf("got %d, want %d", got.Amount, maxUint64)
	}
}

// TestUnsafeMap_float64の精度限界で値がずれる は、float64 が表現できる精度の限界で値がずれることを示す。
// decimal の 12345678.123456789012345678 は float64 では表現できず、近似値に丸められる。
func TestUnsafeMap_float64の精度限界で値がずれる(t *testing.T) {
	// float64 の精度限界により、この値は 12345678.123456789 として表現される
	row := amount.ImportedRow{Amount: 12345678.123456789012345678}
	got := amount.UnsafeMap(row)
	// 小数部が切り捨てられた整数部も、float64 の精度限界で近似値になっている可能性がある
	if got.Amount != 12345678 {
		t.Errorf("got %d, want 12345678", got.Amount)
	}
}
