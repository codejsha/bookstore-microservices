ALTER TABLE stock_history
    ADD CONSTRAINT chk_stock_history_type
        CHECK (change_type IN ('INBOUND', 'OUTBOUND', 'ADJUSTMENT', 'RESERVATION', 'RELEASE'));
ALTER TABLE stock_transfer
    ADD CONSTRAINT chk_transfer_status
        CHECK (status IN ('PENDING', 'IN_TRANSIT', 'COMPLETED', 'CANCELLED'));
ALTER TABLE stock_transfer
    ADD CONSTRAINT chk_transfer_different_warehouses
        CHECK (source_warehouse_id <> target_warehouse_id);
ALTER TABLE stock_audit
    ADD CONSTRAINT chk_audit_status
        CHECK (status IN ('IN_PROGRESS', 'COMPLETED'));
ALTER TABLE monthly_closing
    ADD CONSTRAINT chk_closing_status
        CHECK (status IN ('OPEN', 'CLOSED'));
ALTER TABLE stock_reservation
    ADD CONSTRAINT chk_stock_reservation_status
        CHECK (status IN ('RESERVED', 'RELEASED'));

ALTER TABLE stock
    ADD CONSTRAINT fk_stock__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
ALTER TABLE stock_history
    ADD CONSTRAINT fk_stock_history__stock
        FOREIGN KEY (stock_id) REFERENCES stock (id) ON DELETE CASCADE;
ALTER TABLE stock_transfer
    ADD CONSTRAINT fk_transfer__source_warehouse
        FOREIGN KEY (source_warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
ALTER TABLE stock_transfer
    ADD CONSTRAINT fk_transfer__target_warehouse
        FOREIGN KEY (target_warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
ALTER TABLE stock_audit
    ADD CONSTRAINT fk_audit__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
ALTER TABLE stock_audit_item
    ADD CONSTRAINT fk_audit_item__audit
        FOREIGN KEY (audit_id) REFERENCES stock_audit (id) ON DELETE CASCADE;
ALTER TABLE monthly_closing
    ADD CONSTRAINT fk_closing__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
ALTER TABLE monthly_closing_item
    ADD CONSTRAINT fk_closing_item__closing
        FOREIGN KEY (closing_id) REFERENCES monthly_closing (id) ON DELETE CASCADE;
ALTER TABLE stock_reservation
    ADD CONSTRAINT fk_stock_reservation__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT;
