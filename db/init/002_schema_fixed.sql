-- 修正後のスキーマ
-- BIGINT UNSIGNED を DECIMAL(35, 20) に変更し、
-- amount が必ず正の値であることを CHECK 制約で保証する

-- amount_unsafe は既存カラムのため DEFAULT 0 を付与してテスト時の衝突を回避する
ALTER TABLE transactions
    MODIFY amount_unsafe BIGINT UNSIGNED NOT NULL DEFAULT 0,
    ADD COLUMN amount_fixed DECIMAL(35, 20) NOT NULL DEFAULT 0,
    ADD CONSTRAINT chk_amount_fixed_positive CHECK (amount_fixed >= 0);
