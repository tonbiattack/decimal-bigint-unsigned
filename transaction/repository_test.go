package transaction_test

import (
	"fmt"
	"os"
	"testing"

	appdb "github.com/tonbiattack/decimal-bigint-unsigned/db"
	"github.com/tonbiattack/decimal-bigint-unsigned/transaction"
)

func dsn() string {
	if v := os.Getenv("DSN"); v != "" {
		return v
	}
	return "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
}

func setup(t *testing.T) *transaction.Repository {
	t.Helper()
	conn, err := appdb.Connect(dsn())
	if err != nil {
		t.Skipf("DB接続できないためスキップします: %v", err)
	}
	repo := transaction.NewRepository(conn)

	// テスト前にテーブルをクリア
	if err := conn.Exec("DELETE FROM transactions").Error; err != nil {
		t.Fatalf("テーブルのクリアに失敗しました: %v", err)
	}
	return repo
}

// TestDB_正常な小数値はBIGINT_UNSIGNEDで小数部が消える は、
// float64 の小数値を uint64 にキャストすると小数部が失われることをDBで確認する。
func TestDB_正常な小数値はBIGINT_UNSIGNEDで小数部が消える(t *testing.T) {
	repo := setup(t)

	floatAmount := 12345678.123456789012345678
	tx := transaction.Transaction{
		AmountUnsafe: uint64(floatAmount),
		AmountSafe:   fmt.Sprintf("%.20f", floatAmount),
	}
	if err := repo.Insert(tx); err != nil {
		t.Fatalf("insert に失敗しました: %v", err)
	}

	rows, err := repo.FindAll()
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	got := rows[0]
	t.Logf("amount_unsafe = %d", got.AmountUnsafe)
	t.Logf("amount_safe   = %s", got.AmountSafe)

	// BIGINT UNSIGNED: 小数部が消えて整数部だけになる
	if got.AmountUnsafe != 12345678 {
		t.Errorf("amount_unsafe: got %d, want 12345678", got.AmountUnsafe)
	}
	// DECIMAL: 小数部が保持される（float64 精度まで）
	if got.AmountSafe == "12345678.12345678918063640594" {
		t.Logf("DECIMAL カラムは float64 の近似値で記録される（連携元の精度とは異なる）")
	}
}

// TestDB_極小値はBIGINT_UNSIGNEDでゼロになる は、
// float64 の極小値を uint64 にキャストするとゼロになることをDBで確認する。
func TestDB_極小値はBIGINT_UNSIGNEDでゼロになる(t *testing.T) {
	repo := setup(t)

	floatAmount := 0.000000000000000001
	tx := transaction.Transaction{
		AmountUnsafe: uint64(floatAmount),
		AmountSafe:   fmt.Sprintf("%.20f", floatAmount),
	}
	if err := repo.Insert(tx); err != nil {
		t.Fatalf("insert に失敗しました: %v", err)
	}

	rows, err := repo.FindAll()
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	got := rows[0]
	t.Logf("amount_unsafe = %d", got.AmountUnsafe)
	t.Logf("amount_safe   = %s", got.AmountSafe)

	if got.AmountUnsafe != 0 {
		t.Errorf("amount_unsafe: got %d, want 0", got.AmountUnsafe)
	}
}

// TestDB_負値はBIGINT_UNSIGNEDで最大値に化ける は、
// float64 の負値を uint64 にキャストすると最大値になることをDBで確認する。
// これが今回のトラブルで記録されていた 18446744073709551615 の原因。
func TestDB_負値はBIGINT_UNSIGNEDで最大値に化ける(t *testing.T) {
	repo := setup(t)

	floatAmount := float64(-1)
	tx := transaction.Transaction{
		AmountUnsafe: uint64(floatAmount),
		AmountSafe:   fmt.Sprintf("%.20f", floatAmount),
	}
	if err := repo.Insert(tx); err != nil {
		t.Fatalf("insert に失敗しました: %v", err)
	}

	rows, err := repo.FindAll()
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	got := rows[0]
	t.Logf("amount_unsafe = %d", got.AmountUnsafe)
	t.Logf("amount_safe   = %s", got.AmountSafe)

	// BIGINT UNSIGNED: 最大値に化ける
	want := uint64(18446744073709551615)
	if got.AmountUnsafe != want {
		t.Errorf("amount_unsafe: got %d, want %d", got.AmountUnsafe, want)
	}
	// DECIMAL: 負値のまま保持される
	if got.AmountSafe != "-1.00000000000000000000" {
		t.Errorf("amount_safe: got %s, want -1.00000000000000000000", got.AmountSafe)
	}
}
