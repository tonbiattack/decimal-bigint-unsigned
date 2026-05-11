package main

import (
	"fmt"
	"log"

	appdb "github.com/tonbiattack/decimal-bigint-unsigned/db"
	"github.com/tonbiattack/decimal-bigint-unsigned/transaction"
)

const dsn = "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	conn, err := appdb.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := transaction.NewRepository(conn)

	// DMS 経由で受け取った float64 の値を再現したケース
	cases := []struct {
		label        string
		floatAmount  float64
	}{
		{"正常な小数値", 12345678.123456789012345678},
		{"極小値", 0.000000000000000001},
		{"負値（アンダーフロー）", -1},
	}

	for _, c := range cases {
		tx := transaction.Transaction{
			// 問題のある変換: float64 を uint64 に直接キャスト
			AmountUnsafe: uint64(c.floatAmount),
			// 安全な変換: fmt.Sprintf で文字列化して DECIMAL カラムへ渡す
			AmountSafe: fmt.Sprintf("%.20f", c.floatAmount),
		}
		if err := repo.Insert(tx); err != nil {
			log.Printf("[%s] insert エラー: %v", c.label, err)
			continue
		}
		log.Printf("[%s] insert 完了 — unsafe: %d / safe: %s", c.label, tx.AmountUnsafe, tx.AmountSafe)
	}

	// 登録した結果を取得して表示
	fmt.Println("\n--- DB に記録された値 ---")
	rows, err := repo.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		fmt.Printf("id=%-3d  amount_unsafe=%-25d  amount_safe=%s\n", row.ID, row.AmountUnsafe, row.AmountSafe)
	}
}
