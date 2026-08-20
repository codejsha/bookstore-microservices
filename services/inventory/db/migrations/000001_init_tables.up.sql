CREATE TABLE IF NOT EXISTS warehouse
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID          NOT NULL,
    name       VARCHAR(255)  NOT NULL,
    address    VARCHAR(255),
    capacity   INTEGER       NOT NULL,
    created_at TIMESTAMP(6)  NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT        NOT NULL,
    version    BIGINT        NOT NULL,
    CONSTRAINT uk_warehouse_uid UNIQUE (uid)
);

CREATE INDEX idx_warehouse_name ON warehouse (name);

CREATE TABLE IF NOT EXISTS stock
(
    id            BIGSERIAL PRIMARY KEY,
    uid           UUID         NOT NULL,
    edition_id    BIGINT       NOT NULL,
    edition_uid   UUID         NOT NULL,
    warehouse_id  BIGINT       NOT NULL,
    warehouse_uid UUID         NOT NULL,
    quantity      INTEGER      NOT NULL,
    created_at    TIMESTAMP(6) NOT NULL,
    updated_at    TIMESTAMP(6),
    deleted_at    TIMESTAMP(6),
    actor         BIGINT       NOT NULL,
    version       BIGINT       NOT NULL,
    CONSTRAINT fk_stock__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT uk_stock_uid UNIQUE (uid),
    CONSTRAINT uq_stock_edition_warehouse UNIQUE (edition_id, warehouse_id)
);

CREATE INDEX idx_stock_edition_id ON stock (edition_id);
CREATE INDEX idx_stock_warehouse_id ON stock (warehouse_id);

CREATE TABLE IF NOT EXISTS stock_history
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    stock_id    BIGINT       NOT NULL,
    stock_uid   UUID         NOT NULL,
    change_type VARCHAR(50)  NOT NULL,
    reason      VARCHAR(255),
    change_qty  INTEGER      NOT NULL,
    created_at  TIMESTAMP(6) NOT NULL,
    updated_at  TIMESTAMP(6),
    deleted_at  TIMESTAMP(6),
    actor       BIGINT       NOT NULL,
    version     BIGINT       NOT NULL,
    CONSTRAINT fk_stock_history__stock
        FOREIGN KEY (stock_id) REFERENCES stock (id) ON DELETE CASCADE,
    CONSTRAINT chk_stock_history_type
        CHECK (change_type IN ('INBOUND', 'OUTBOUND', 'ADJUSTMENT', 'RESERVATION', 'RELEASE')),
    CONSTRAINT uk_stock_history_uid UNIQUE (uid)
);

CREATE INDEX idx_stock_history_stock_id ON stock_history (stock_id);

CREATE TABLE IF NOT EXISTS stock_transfer
(
    id                   BIGSERIAL PRIMARY KEY,
    uid                  UUID         NOT NULL,
    edition_id           BIGINT       NOT NULL,
    edition_uid          UUID         NOT NULL,
    source_warehouse_id  BIGINT       NOT NULL,
    source_warehouse_uid UUID         NOT NULL,
    target_warehouse_id  BIGINT       NOT NULL,
    target_warehouse_uid UUID         NOT NULL,
    quantity             INTEGER      NOT NULL,
    status               VARCHAR(20)  NOT NULL,
    reason               VARCHAR(255),
    completed_at         TIMESTAMP(6),
    created_at           TIMESTAMP(6) NOT NULL,
    updated_at           TIMESTAMP(6),
    deleted_at           TIMESTAMP(6),
    actor                BIGINT       NOT NULL,
    version              BIGINT       NOT NULL,
    CONSTRAINT fk_transfer__source_warehouse
        FOREIGN KEY (source_warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT fk_transfer__target_warehouse
        FOREIGN KEY (target_warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT chk_transfer_status
        CHECK (status IN ('PENDING', 'IN_TRANSIT', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT chk_transfer_different_warehouses
        CHECK (source_warehouse_id <> target_warehouse_id),
    CONSTRAINT uk_transfer_uid UNIQUE (uid)
);

CREATE INDEX idx_transfer_edition_id ON stock_transfer (edition_id);
CREATE INDEX idx_transfer_source ON stock_transfer (source_warehouse_id);
CREATE INDEX idx_transfer_target ON stock_transfer (target_warehouse_id);
CREATE INDEX idx_transfer_status ON stock_transfer (status);

CREATE TABLE IF NOT EXISTS stock_audit
(
    id            BIGSERIAL PRIMARY KEY,
    uid           UUID         NOT NULL,
    warehouse_id  BIGINT       NOT NULL,
    warehouse_uid UUID         NOT NULL,
    status        VARCHAR(20)  NOT NULL,
    notes         VARCHAR(500),
    completed_at  TIMESTAMP(6),
    created_at    TIMESTAMP(6) NOT NULL,
    updated_at    TIMESTAMP(6),
    deleted_at    TIMESTAMP(6),
    actor         BIGINT       NOT NULL,
    version       BIGINT       NOT NULL,
    CONSTRAINT fk_audit__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT chk_audit_status
        CHECK (status IN ('IN_PROGRESS', 'COMPLETED')),
    CONSTRAINT uk_audit_uid UNIQUE (uid)
);

CREATE INDEX idx_audit_warehouse_id ON stock_audit (warehouse_id);
CREATE INDEX idx_audit_status ON stock_audit (status);

CREATE TABLE IF NOT EXISTS stock_audit_item
(
    id              BIGSERIAL PRIMARY KEY,
    audit_id        BIGINT  NOT NULL,
    edition_id      BIGINT  NOT NULL,
    edition_uid     UUID    NOT NULL,
    system_quantity INTEGER NOT NULL,
    actual_quantity INTEGER NOT NULL,
    difference      INTEGER NOT NULL,
    CONSTRAINT fk_audit_item__audit
        FOREIGN KEY (audit_id) REFERENCES stock_audit (id) ON DELETE CASCADE
);

CREATE INDEX idx_audit_item_audit_id ON stock_audit_item (audit_id);

CREATE TABLE IF NOT EXISTS monthly_closing
(
    id            BIGSERIAL PRIMARY KEY,
    uid           UUID         NOT NULL,
    warehouse_id  BIGINT       NOT NULL,
    warehouse_uid UUID         NOT NULL,
    year          INTEGER      NOT NULL,
    month         INTEGER      NOT NULL,
    status        VARCHAR(20)  NOT NULL,
    closed_at     TIMESTAMP(6),
    created_at    TIMESTAMP(6) NOT NULL,
    updated_at    TIMESTAMP(6),
    deleted_at    TIMESTAMP(6),
    actor         BIGINT       NOT NULL,
    version       BIGINT       NOT NULL,
    CONSTRAINT fk_closing__warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE RESTRICT,
    CONSTRAINT chk_closing_status
        CHECK (status IN ('OPEN', 'CLOSED')),
    CONSTRAINT uk_closing_uid UNIQUE (uid),
    CONSTRAINT uq_closing_warehouse_period UNIQUE (warehouse_id, year, month)
);

CREATE INDEX idx_closing_warehouse_id ON monthly_closing (warehouse_id);

CREATE TABLE IF NOT EXISTS monthly_closing_item
(
    id                BIGSERIAL PRIMARY KEY,
    closing_id        BIGINT  NOT NULL,
    edition_id        BIGINT  NOT NULL,
    edition_uid       UUID    NOT NULL,
    opening_quantity  INTEGER NOT NULL,
    inbound_quantity  INTEGER NOT NULL,
    outbound_quantity INTEGER NOT NULL,
    adjust_quantity   INTEGER NOT NULL,
    closing_quantity  INTEGER NOT NULL,
    CONSTRAINT fk_closing_item__closing
        FOREIGN KEY (closing_id) REFERENCES monthly_closing (id) ON DELETE CASCADE
);

CREATE INDEX idx_closing_item_closing_id ON monthly_closing_item (closing_id);
