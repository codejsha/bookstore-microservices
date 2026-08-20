ALTER TABLE stock_reservation
    DROP CONSTRAINT uk_stock_reservation_key_warehouse;

ALTER TABLE stock_reservation
    ADD CONSTRAINT uk_stock_reservation_key UNIQUE (reservation_key);
