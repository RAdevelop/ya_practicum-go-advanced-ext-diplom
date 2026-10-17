CREATE TABLE IF NOT EXISTS balances (
    customer_id BIGINT PRIMARY KEY NOT NULL CHECK (customer_id >= 0),
    "current" NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK ("current" >= 0),
    withdrawn NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (withdrawn >= 0),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

COMMENT ON TABLE  balances               IS 'текущий баланс покупателей';
COMMENT ON COLUMN balances.customer_id   IS 'id покупателя, которому принадлежит баланс';
COMMENT ON COLUMN balances.current       IS 'текущее количество накопленных баллов доступных для списания';
COMMENT ON COLUMN balances.withdrawn     IS 'текущее количество списанных баллов';
COMMENT ON COLUMN balances.updated_at    IS 'когда баланс был изменен';