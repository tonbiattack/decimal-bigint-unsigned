package amount

// ImportedRow は DMS 経由でインポートされたデータの構造体。
// DMS が decimal カラムを float64 に変換するため、フィールドが float64 になっている。
type ImportedRow struct {
	Amount float64
}

// DestinationRow はインポート先システムの構造体。
// 金額を uint64 で持つ設計になっている。
type DestinationRow struct {
	Amount uint64
}

// UnsafeMapFromDMS は DMS インポート行をインポート先の struct に変換する。
// float64 を uint64 にキャストするため、精度欠損・アンダーフローが発生する危険な実装例。
func UnsafeMapFromDMS(row ImportedRow) DestinationRow {
	return DestinationRow{
		Amount: uint64(row.Amount),
	}
}
