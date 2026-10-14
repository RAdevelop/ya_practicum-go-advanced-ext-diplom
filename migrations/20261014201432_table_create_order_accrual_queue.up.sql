CREATE TABLE IF NOT EXISTS order_accrual_queue (
    order_number VARCHAR(20) PRIMARY KEY NOT NULL CHECK (order_number != ''),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    next_check_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    attempts SMALLINT NOT NULL DEFAULT 0 CHECK (attempts >= 0)
);

CREATE INDEX IF NOT EXISTS idx_order_accrual_queue_next_check_at ON order_accrual_queue (next_check_at DESC);
CREATE INDEX IF NOT EXISTS idx_order_accrual_queue_created_at ON order_accrual_queue (created_at ASC);

COMMENT ON TABLE  order_accrual_queue                  IS 'очередь заказов для получения начислений';
COMMENT ON COLUMN order_accrual_queue.order_number     IS 'уникальный номер заказа';
COMMENT ON COLUMN order_accrual_queue.created_at       IS 'когда заказ был создан';
COMMENT ON COLUMN order_accrual_queue.updated_at       IS 'когда заказ был изменен';
COMMENT ON COLUMN order_accrual_queue.next_check_at    IS 'когда делать следующий запрос для начислений';
COMMENT ON COLUMN order_accrual_queue.attempts         IS 'количество сделанных попыток получения начисления';