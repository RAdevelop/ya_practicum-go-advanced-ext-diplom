CREATE TABLE IF NOT EXISTS withdrawals (
    order_id BIGINT PRIMARY KEY NOT NULL CHECK (order_id >= 0),
    "sum" NUMERIC(10, 2) NOT NULL CHECK ("sum" >= 0),
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_withdrawals_processed_at ON withdrawals (processed_at DESC);

COMMENT ON TABLE  withdrawals               IS 'информация о выводе средств';
COMMENT ON COLUMN withdrawals.order_id      IS 'id заказа';
COMMENT ON COLUMN withdrawals.sum           IS 'сумма списания';
COMMENT ON COLUMN withdrawals.processed_at  IS 'когда было списание';