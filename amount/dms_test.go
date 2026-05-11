package amount_test

import (
	"testing"

	"github.com/tonbiattack/decimal-bigint-unsigned/amount"
)

// このファイルは今回のトラブルを再現するテストです。
//
// 問題の経路:
//   連携元 decimal → DMS（float64 に変換）→ Go（uint64 にキャスト）→ BIGINT UNSIGNED
//
// 発生した事象:
//   1. 18446744073709551615（BIGINT UNSIGNED の最大値）が記録された
//   2. 小数点以下が消失した

// TestTrouble_最大値がDBに記録された再現 は、float64 の負値を uint64 に直接キャストすると
// 最大値になることを示す。
// Go の仕様では float64 を整数型にキャストするとき、値が変換先の範囲外なら結果は実装依存。
// 負値の float64 を uint64 に直接キャストすると、このGoバージョンでは最大値になる。
func TestTrouble_最大値がDBに記録された再現(t *testing.T) {
	// DMS が連携元の負値 decimal を float64 として渡してきた状態を再現
	row := amount.ImportedRow{Amount: -1}
	got := amount.UnsafeMapFromDMS(row)

	// 期待値: 本来の取引金額（例: 1）
	// 実際の値: 18446744073709551615（BIGINT UNSIGNED の最大値）
	want := uint64(18446744073709551615)
	if got.Amount != want {
		t.Errorf("got %d, want %d", got.Amount, want)
	}
	t.Logf("float64(-1) を uint64 に直接キャストすると: %d", got.Amount)
}

// TestTrouble_小数点以下が消失した再現 は、DMS が float64 に変換した時点で
// decimal の精度が失われ、さらに uint64 キャストで小数部が消えることを示す。
func TestTrouble_小数点以下が消失した再現(t *testing.T) {
	// DMS が連携元の高精度 decimal を float64 として渡してきた状態を再現
	// 連携元の実際の値: 12345678.123456789012345678
	// DMS 経由で float64 になった値: 12345678.123456789（精度がすでに失われている）
	row := amount.ImportedRow{Amount: 12345678.123456789012345678}
	got := amount.UnsafeMapFromDMS(row)

	// 期待値: 12345678.123456789012345678（小数部を含む）
	// 実際の値: 12345678（小数部が消えている）
	if got.Amount != 12345678 {
		t.Errorf("got %d, want 12345678", got.Amount)
	}
	t.Logf("float64(12345678.123456789012345678) を uint64 にキャストすると: %d", got.Amount)
}

// TestTrouble_極小値がゼロになった再現 は、0.000000000000000001 のような極小値が
// float64 → uint64 の変換でゼロになることを示す。
func TestTrouble_極小値がゼロになった再現(t *testing.T) {
	row := amount.ImportedRow{Amount: 0.000000000000000001}
	got := amount.UnsafeMapFromDMS(row)

	if got.Amount != 0 {
		t.Errorf("got %d, want 0", got.Amount)
	}
	t.Logf("float64(0.000000000000000001) を uint64 にキャストすると: %d", got.Amount)
}

// TestTrouble_float64時点ですでに精度が失われている は、DMS による float64 変換の時点で
// すでに decimal の精度が失われていることを示す。uint64 キャスト以前の問題。
func TestTrouble_float64時点ですでに精度が失われている(t *testing.T) {
	// 連携元の値
	original := "12345678.123456789012345678"

	// DMS が float64 に変換した結果（Goのfloat64リテラルで再現）
	asFloat := float64(12345678.123456789012345678)

	// float64 で表現できる値を文字列に戻すと、元の精度には戻らない
	d, err := amount.SafeMapFromDMS(amount.ImportedRow{Amount: asFloat})
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}

	if d.String() == original {
		t.Errorf("精度が保たれてしまっています: got %s", d.String())
	}
	t.Logf("連携元の値:         %s", original)
	t.Logf("float64 経由の値:   %s", d.String())
}
