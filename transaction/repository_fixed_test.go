package transaction_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/tonbiattack/decimal-bigint-unsigned/amount"
	appdb "github.com/tonbiattack/decimal-bigint-unsigned/db"
	"github.com/tonbiattack/decimal-bigint-unsigned/transaction"
)

func setupFixed(t *testing.T) *transaction.RepositoryFixed {
	t.Helper()
	conn, err := appdb.Connect(dsn())
	if err != nil {
		t.Skipf("DB接続できないためスキップします: %v", err)
	}
	repo := transaction.NewRepositoryFixed(conn)
	if err := conn.Exec("DELETE FROM transactions").Error; err != nil {
		t.Fatalf("テーブルのクリアに失敗しました: %v", err)
	}
	return repo
}

// TestDBFixed_正常な小数値は精度を保ったまま記録される は、
// 文字列 → decimal → DECIMAL(35,20) の経路では精度が失われないことをDBで確認する。
func TestDBFixed_正常な小数値は精度を保ったまま記録される(t *testing.T) {
	repo := setupFixed(t)

	input := "12345678.123456789012345678"
	d, err := amount.SafeMapFromDMSFixed(amount.ImportedRowFixed{Amount: input})
	if err != nil {
		t.Fatalf("変換エラー: %v", err)
	}

	if err := repo.Insert(transaction.TransactionFixed{AmountFixed: d}); err != nil {
		t.Fatalf("insert に失敗しました: %v", err)
	}

	rows, err := repo.FindAll()
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	got := rows[0].AmountFixed
	want := decimal.RequireFromString(input)
	t.Logf("amount_fixed = %s", got.String())

	if !got.Equal(want) {
		t.Errorf("got %s, want %s", got.String(), want.String())
	}
}

// TestDBFixed_極小値は精度を保ったまま記録される は、
// 極小値が DECIMAL カラムでゼロにならずに保持されることをDBで確認する。
func TestDBFixed_極小値は精度を保ったまま記録される(t *testing.T) {
	repo := setupFixed(t)

	input := "0.000000000000000001"
	d, err := amount.SafeMapFromDMSFixed(amount.ImportedRowFixed{Amount: input})
	if err != nil {
		t.Fatalf("変換エラー: %v", err)
	}

	if err := repo.Insert(transaction.TransactionFixed{AmountFixed: d}); err != nil {
		t.Fatalf("insert に失敗しました: %v", err)
	}

	rows, err := repo.FindAll()
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	got := rows[0].AmountFixed
	want := decimal.RequireFromString(input)
	t.Logf("amount_fixed = %s", got.String())

	if !got.Equal(want) {
		t.Errorf("got %s, want %s", got.String(), want.String())
	}
}

// TestDBFixed_負値は変換時点でエラーになる は、
// 負値が SafeMapFromDMSFixed でエラーになり DB に到達しないことを確認する。
func TestDBFixed_負値は変換時点でエラーになる(t *testing.T) {
	_, err := amount.SafeMapFromDMSFixed(amount.ImportedRowFixed{Amount: "-1"})
	if err == nil {
		t.Error("負値はエラーになるべきですが、nil でした")
	}
	t.Logf("負値の変換エラー: %v", err)
}
