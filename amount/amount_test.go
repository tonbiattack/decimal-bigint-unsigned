package amount_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/tonbiattack/decimal-bigint-unsigned/amount"
)

// uint64 の最大値
const maxUint64 = uint64(18446744073709551615)

// TestUnsafeConvertToUint64_正の小数は整数部だけになる は、正の小数が uint64 変換で小数部を失うことを示す。
func TestUnsafeConvertToUint64_正の小数は整数部だけになる(t *testing.T) {
	d := decimal.RequireFromString("12345678.123456789012345678")
	got := amount.UnsafeConvertToUint64(d)
	// 小数部が切り捨てられる
	if got != 12345678 {
		t.Errorf("got %d, want 12345678", got)
	}
}

// TestUnsafeConvertToUint64_負数はuint64最大値に化ける は、負数を uint64 に変換すると最大値になることを示す。
func TestUnsafeConvertToUint64_負数はuint64最大値に化ける(t *testing.T) {
	d := decimal.RequireFromString("-1")
	got := amount.UnsafeConvertToUint64(d)
	if got != maxUint64 {
		t.Errorf("got %d, want %d", got, maxUint64)
	}
}

// TestUnsafeConvertToUint64_非常に小さい正の小数はゼロになる は、0.000...001 のような値がゼロに化けることを示す。
func TestUnsafeConvertToUint64_非常に小さい正の小数はゼロになる(t *testing.T) {
	d := decimal.RequireFromString("0.000000000000000001")
	got := amount.UnsafeConvertToUint64(d)
	if got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

// TestUnsafeConvertDecimalDirectlyToUint64_正の小数は整数部だけになる は、小数部が切り捨てられることを示す。
func TestUnsafeConvertDecimalDirectlyToUint64_正の小数は整数部だけになる(t *testing.T) {
	d := decimal.RequireFromString("12345678.123456789012345678")
	got := amount.UnsafeConvertDecimalDirectlyToUint64(d)
	if got != 12345678 {
		t.Errorf("got %d, want 12345678", got)
	}
}

// TestUnsafeConvertDecimalDirectlyToUint64_負数は絶対値になる は、float64 経由と異なり最大値ではなく絶対値になることを示す。
func TestUnsafeConvertDecimalDirectlyToUint64_負数は絶対値になる(t *testing.T) {
	d := decimal.RequireFromString("-1")
	got := amount.UnsafeConvertDecimalDirectlyToUint64(d)
	// float64 経由（uint64 最大値）とは異なり、絶対値の 1 になる
	if got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

// TestUnsafeConvertDecimalDirectlyToUint64_負の小数は絶対値の整数部になる は、負の小数が絶対値の整数部に化けることを示す。
func TestUnsafeConvertDecimalDirectlyToUint64_負の小数は絶対値の整数部になる(t *testing.T) {
	d := decimal.RequireFromString("-12345678.5")
	got := amount.UnsafeConvertDecimalDirectlyToUint64(d)
	if got != 12345678 {
		t.Errorf("got %d, want 12345678", got)
	}
}

// TestUnsafeConvertDecimalDirectlyToUint64_非常に小さい正の小数はゼロになる は、極小値がゼロに化けることを示す。
func TestUnsafeConvertDecimalDirectlyToUint64_非常に小さい正の小数はゼロになる(t *testing.T) {
	d := decimal.RequireFromString("0.000000000000000001")
	got := amount.UnsafeConvertDecimalDirectlyToUint64(d)
	if got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

// TestSafeConvertToString_精度を保ったまま文字列になる は、decimal を文字列変換すると精度が保たれることを示す。
func TestSafeConvertToString_精度を保ったまま文字列になる(t *testing.T) {
	input := "12345678.123456789012345678"
	d := decimal.RequireFromString(input)
	got := amount.SafeConvertToString(d)
	if got != input {
		t.Errorf("got %s, want %s", got, input)
	}
}

// TestSafeConvertToString_非常に小さい値も精度が保たれる は、極小値でも文字列変換では精度が保たれることを示す。
func TestSafeConvertToString_非常に小さい値も精度が保たれる(t *testing.T) {
	input := "0.000000000000000001"
	d := decimal.RequireFromString(input)
	got := amount.SafeConvertToString(d)
	if got != input {
		t.Errorf("got %s, want %s", got, input)
	}
}

// TestSafeConvertToDecimal_文字列からdecimalに変換できる は、文字列から decimal に変換できることを示す。
func TestSafeConvertToDecimal_文字列からdecimalに変換できる(t *testing.T) {
	input := "12345678.123456789012345678"
	d, err := amount.SafeConvertToDecimal(input)
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if d.String() != input {
		t.Errorf("got %s, want %s", d.String(), input)
	}
}

// TestSafeConvertToDecimal_不正な文字列はエラーになる は、不正な入力に対してエラーが返ることを示す。
func TestSafeConvertToDecimal_不正な文字列はエラーになる(t *testing.T) {
	_, err := amount.SafeConvertToDecimal("not-a-number")
	if err == nil {
		t.Error("エラーが返るべきですが、nil でした")
	}
}

// TestUnsafeMapFromDMS_小数部が切り捨てられる は、DMS が float64 で渡した小数値が uint64 変換で切り捨てられることを示す。
func TestUnsafeMapFromDMS_小数部が切り捨てられる(t *testing.T) {
	// DB の decimal カラムが DMS 経由で float64 になった値
	row := amount.ImportedRow{Amount: 12345678.123456789}
	got := amount.UnsafeMapFromDMS(row)
	if got.Amount != 12345678 {
		t.Errorf("got %d, want 12345678", got.Amount)
	}
}

// TestUnsafeMapFromDMS_極小値がゼロになる は、0.000...001 のような値がゼロに化けることを示す。
func TestUnsafeMapFromDMS_極小値がゼロになる(t *testing.T) {
	row := amount.ImportedRow{Amount: 0.000000000000000001}
	got := amount.UnsafeMapFromDMS(row)
	if got.Amount != 0 {
		t.Errorf("got %d, want 0", got.Amount)
	}
}

// TestUnsafeMapFromDMS_負数がuint64最大値に化ける は、DMS が負数を float64 で渡したとき uint64 への変換で最大値になることを示す。
// Go の仕様では符号付き整数から符号なし整数への変換は 2^N を法とした値になるため、
// int64(-1) を uint64 に変換すると 18446744073709551615 になる。
func TestUnsafeMapFromDMS_負数がuint64最大値に化ける(t *testing.T) {
	row := amount.ImportedRow{Amount: -1}
	got := amount.UnsafeMapFromDMS(row)
	if got.Amount != maxUint64 {
		t.Errorf("got %d, want %d", got.Amount, maxUint64)
	}
}
