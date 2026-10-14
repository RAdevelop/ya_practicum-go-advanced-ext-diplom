CREATE TABLE IF NOT EXISTS withdrawals (
    order_number VARCHAR(20) PRIMARY KEY NOT NULL CHECK (order_number != ''),
    customer_id BIGINT NOT NULL CHECK (customer_id >= 0),
    "sum" NUMERIC(10, 2) NOT NULL CHECK ("sum" >= 0),
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_withdrawals_customer_processed_at ON withdrawals (customer_id, processed_at DESC);

COMMENT ON TABLE  withdrawals               IS 'информация о выводе средств';
COMMENT ON COLUMN withdrawals.order_number  IS 'номер заказа';
COMMENT ON COLUMN withdrawals.customer_id   IS 'уникальный идентификатор покупателя';
COMMENT ON COLUMN withdrawals.sum           IS 'сумма списания';
COMMENT ON COLUMN withdrawals.processed_at  IS 'когда было списание';