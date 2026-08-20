ALTER TABLE stock_reservation
    DROP CONSTRAINT uk_stock_reservation_key;

ALTER TABLE stock_reservation
    ADD CONSTRAINT uk_stock_reservation_key_warehouse UNIQUE (reservation_key, warehouse_id);
