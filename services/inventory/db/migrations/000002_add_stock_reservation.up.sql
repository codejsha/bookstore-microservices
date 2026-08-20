CREATE TABLE IF NOT EXISTS stock_reservation
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    reservation_key VARCHAR(255) NOT NULL,
    order_ref       VARCHAR(255) NOT NULL,
    edition_id      BIGINT       NOT NULL,
    edition_uid     UUID         NOT NULL,
    warehouse_id    BIGINT       NOT NULL,
    warehouse_uid   UUID         NOT NULL,
    quantity        INTEGER      NOT NULL,
    status          VARCHAR(20)  NOT NULL,
    released_at     TIMESTAMP(6),
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    actor           BIGINT       NOT NULL,
    version         BIGINT       NOT NULL,
    CONSTRAINT fk_stock_reservation__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT chk_stock_reservation_status
        CHECK (status IN ('RESERVED', 'RELEASED')),
    CONSTRAINT uk_stock_reservation_uid UNIQUE (uid),
    CONSTRAINT uk_stock_reservation_key UNIQUE (reservation_key)
);

CREATE INDEX idx_stock_reservation_order_ref ON stock_reservation (order_ref);
CREATE INDEX idx_stock_reservation_edition_id ON stock_reservation (edition_id);
CREATE INDEX idx_stock_reservation_warehouse_id ON stock_reservation (warehouse_id);
