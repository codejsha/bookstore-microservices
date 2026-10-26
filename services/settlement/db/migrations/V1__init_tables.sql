-- ------------------------------------------------------------
-- 1. Daily settlement buckets
--    One row per (settlement_date, currency, payment_method) close.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS daily_settlement
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    settlement_date DATE         NOT NULL,
    currency        CHAR(3)      NOT NULL,
    payment_method  VARCHAR(32),
    gross_amount    BIGINT       NOT NULL,
    refund_amount   BIGINT       NOT NULL,
    fee_amount      BIGINT       NOT NULL,
    net_amount      BIGINT       NOT NULL,
    payment_count   INT          NOT NULL,
    refund_count    INT          NOT NULL,
    status          VARCHAR(20)  NOT NULL,
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    CONSTRAINT uk_daily_settlement_uid UNIQUE (uid),
    CONSTRAINT uk_settlement_bucket UNIQUE (settlement_date, currency, payment_method)
);

CREATE INDEX idx_settlement_date_status ON daily_settlement (settlement_date, status);


-- ------------------------------------------------------------
-- 2. Settlement detail lines
--    Individual payment/refund lines feeding a bucket. Source data lives in the
--    payment service's DB, so there are no cross-service foreign keys here.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_detail
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    settlement_date DATE         NOT NULL,
    source_type     VARCHAR(10)  NOT NULL,
    source_id       VARCHAR(64)  NOT NULL,
    payment_id      VARCHAR(64)  NOT NULL,
    amount          BIGINT       NOT NULL,
    currency        CHAR(3)      NOT NULL,
    payment_method  VARCHAR(32),
    occurred_at     TIMESTAMP(6) NOT NULL,
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    CONSTRAINT uk_settlement_detail_uid UNIQUE (uid),
    CONSTRAINT uk_detail_source UNIQUE (source_type, source_id)
);

CREATE INDEX idx_settlement_detail_date ON settlement_detail (settlement_date);


-- ------------------------------------------------------------
-- 3. Settlement job run log
--    Audit trail of settlement batch executions (Spring Batch).
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_job_run
(
    id               BIGSERIAL PRIMARY KEY,
    uid              UUID          NOT NULL,
    target_date      DATE          NOT NULL,
    job_execution_id BIGINT,
    status           VARCHAR(20)   NOT NULL,
    started_at       TIMESTAMP(6)  NOT NULL,
    finished_at      TIMESTAMP(6),
    message          VARCHAR(1024),
    created_at       TIMESTAMP(6)  NOT NULL,
    updated_at       TIMESTAMP(6),
    deleted_at       TIMESTAMP(6),
    CONSTRAINT uk_settlement_job_run_uid UNIQUE (uid)
);

CREATE INDEX idx_settlement_job_run_target ON settlement_job_run (target_date);
