package main

import (
	"fmt"
	"log"

	"github.com/tonbiattack/decimal-bigint-unsigned/amount"
	appdb "github.com/tonbiattack/decimal-bigint-unsigned/db"
	"github.com/tonbiattack/decimal-bigint-unsigned/transaction"
)

const dsn = "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	conn, err := appdb.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := transaction.NewRepositoryFixed(conn)

	// DMS のマッピングを文字列に変更した想定のケース
	cases := []struct {
		label  string
		amount string
	}{
		{"正常な小数値", "12345678.123456789012345678"},
		{"極小値", "0.000000000000000001"},
		{"負値（エラーになるべきケース）", "-1"},
	}

	for _, c := range cases {
		d, err := amount.SafeMapFromDMSFixed(amount.ImportedRowFixed{Amount: c.amount})
		if err != nil {
			log.Printf("[%s] 変換エラー（上流の異常を検知）: %v", c.label, err)
			continue
		}

		tx := transaction.TransactionFixed{AmountFixed: d}
		if err := repo.Insert(tx); err != nil {
			log.Printf("[%s] insert エラー: %v", c.label, err)
			continue
		}
		log.Printf("[%s] insert 完了 — amount: %s", c.label, d.String())
	}

	fmt.Println("\n--- DB に記録された値 ---")
	rows, err := repo.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		fmt.Printf("id=%-3d  amount_fixed=%s\n", row.ID, row.AmountFixed.String())
	}
}
