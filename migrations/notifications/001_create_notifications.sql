CREATE TABLE notifications (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  event_type TEXT NOT NULL CHECK (event_type IN ('ORDER_CREATED','ORDER_UPDATED','ORDER_DELETED')),
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_notifications_order_id_created_at ON notifications(order_id,created_at DESC);
