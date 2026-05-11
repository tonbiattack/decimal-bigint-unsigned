# decimal-bigint-unsigned

`decimal` を `float64` や `uint64` に変換することで発生する精度欠損・アンダーフローのサンプルコードです。

記事: [decimal を BIGINT UNSIGNED / uint64 で扱ってはいけない理由]()

## 構成

```
amount/
├── amount.go       # 危険な変換と安全な変換の実装例
└── amount_test.go  # 各変換の挙動を示すテスト
```

## 何を示しているか

### 危険な変換（UnsafeConvertToUint64）

`decimal` を `float64` → `uint64` と経由して変換します。次の問題が起きます。

- 小数部が切り捨てられる
- 負数を変換すると `uint64` の最大値（`18446744073709551615`）に化ける
- 極小値（`0.000000000000000001` など）がゼロになる

### 安全な変換（SafeConvertToString / SafeConvertToDecimal）

`shopspring/decimal` を使い、文字列を介して精度を保ったまま扱います。

## テスト実行

```bash
go test ./...
```
