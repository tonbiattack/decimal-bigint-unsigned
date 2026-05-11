CREATE TABLE transactions (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    -- 問題のある型: 整数しか保持できず、範囲外の値は実装依存になる
    amount_unsafe  BIGINT UNSIGNED NOT NULL,
    -- 正しい型: 小数18桁まで保持できる
    amount_safe    DECIMAL(35, 20) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
