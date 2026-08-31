ALTER TABLE stock_history DROP CONSTRAINT IF EXISTS chk_stock_history_type;
ALTER TABLE stock_transfer DROP CONSTRAINT IF EXISTS chk_transfer_status;
ALTER TABLE stock_transfer DROP CONSTRAINT IF EXISTS chk_transfer_different_warehouses;
ALTER TABLE stock_audit DROP CONSTRAINT IF EXISTS chk_audit_status;
ALTER TABLE monthly_closing DROP CONSTRAINT IF EXISTS chk_closing_status;
ALTER TABLE stock_reservation DROP CONSTRAINT IF EXISTS chk_stock_reservation_status;

ALTER TABLE stock DROP CONSTRAINT IF EXISTS fk_stock__warehouse;
ALTER TABLE stock_history DROP CONSTRAINT IF EXISTS fk_stock_history__stock;
ALTER TABLE stock_transfer DROP CONSTRAINT IF EXISTS fk_transfer__source_warehouse;
ALTER TABLE stock_transfer DROP CONSTRAINT IF EXISTS fk_transfer__target_warehouse;
ALTER TABLE stock_audit DROP CONSTRAINT IF EXISTS fk_audit__warehouse;
ALTER TABLE stock_audit_item DROP CONSTRAINT IF EXISTS fk_audit_item__audit;
ALTER TABLE monthly_closing DROP CONSTRAINT IF EXISTS fk_closing__warehouse;
ALTER TABLE monthly_closing_item DROP CONSTRAINT IF EXISTS fk_closing_item__closing;
ALTER TABLE stock_reservation DROP CONSTRAINT IF EXISTS fk_stock_reservation__warehouse;
