CREATE TABLE IF NOT EXISTS order_accrual_queue (
    order_number VARCHAR(20) PRIMARY KEY NOT NULL CHECK (order_number != ''),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    next_check_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    attempts SMALLINT NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_order_accrual_queue_next_check_at ON order_accrual_queue (next_check_at);

COMMENT ON TABLE  order_accrual_queue                  IS 'очередь заказов для получения начислений';
COMMENT ON COLUMN order_accrual_queue.order_number     IS 'уникальный номер заказа';
COMMENT ON COLUMN order_accrual_queue.updated_at       IS 'когда заказ был изменен';
COMMENT ON COLUMN order_accrual_queue.next_check_at    IS 'когда делать следующий запрос для начислений';