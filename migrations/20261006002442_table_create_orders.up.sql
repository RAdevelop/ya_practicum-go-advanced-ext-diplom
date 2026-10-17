CREATE TABLE IF NOT EXISTS orders (
    "number" VARCHAR(20) PRIMARY KEY NOT NULL CHECK ("number" != ''),
    customer_id BIGINT NOT NULL CHECK (customer_id >= 0),
    status VARCHAR(10) NOT NULL CHECK (status != ''),
    accrual NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (accrual >= 0),
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_uploaded ON orders (customer_id, uploaded_at DESC);

COMMENT ON TABLE  orders               IS 'заказы покупателей системы лояльности';
COMMENT ON COLUMN orders.number        IS 'уникальный номер заказа';
COMMENT ON COLUMN orders.customer_id   IS 'id покупателя, которому принадлежит заказ';
COMMENT ON COLUMN orders.status        IS 'статус обработки расчёта';
COMMENT ON COLUMN orders.accrual       IS 'начисленные баллы';
COMMENT ON COLUMN orders.uploaded_at   IS 'когда заказ был загружен';
COMMENT ON COLUMN orders.updated_at    IS 'когда заказ был изменен';