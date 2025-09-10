SET NAMES utf8mb4;

-- ------------------------------------------------------------
-- 1. Daily settlement buckets
--    One row per (settlement_date, currency, payment_method) close.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS daily_settlement (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid             BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    settlement_date DATE         NOT NULL COMMENT 'Business date being settled',
    currency        CHAR(3)      NOT NULL COMMENT 'ISO 4217 currency code (e.g. USD, KRW)',
    payment_method  VARCHAR(32)  NULL     COMMENT 'Payment method bucket (card, wallet, ...); NULL = all methods aggregated',
    gross_amount    BIGINT       NOT NULL COMMENT 'Captured payment total (smallest currency unit)',
    refund_amount   BIGINT       NOT NULL COMMENT 'Refunded total (smallest currency unit)',
    fee_amount      BIGINT       NOT NULL COMMENT 'PG/processing fee total (smallest currency unit)',
    net_amount      BIGINT       NOT NULL COMMENT 'Net settled = gross - refund - fee (smallest currency unit)',
    payment_count   INT          NOT NULL COMMENT 'Number of captured payments in this bucket',
    refund_count    INT          NOT NULL COMMENT 'Number of refunds in this bucket',
    status          VARCHAR(20)  NOT NULL COMMENT 'OPEN, CONFIRMED, DISCREPANCY',
    created_at      DATETIME(6)  NOT NULL,
    updated_at      DATETIME(6)  NULL,
    deleted_at      DATETIME(6)  NULL,

    UNIQUE KEY uk_daily_settlement_uid (uid),
    UNIQUE KEY uk_settlement_bucket (settlement_date, currency, payment_method),
    INDEX idx_settlement_date_status (settlement_date, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Daily sales-close / PG reconciliation buckets';


-- ------------------------------------------------------------
-- 2. Settlement detail lines
--    Individual payment/refund lines feeding a bucket. Source data lives in the
--    payment service's DB, so there are no cross-service foreign keys here.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_detail (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid             BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    settlement_date DATE         NOT NULL COMMENT 'Business date this line belongs to',
    source_type     VARCHAR(10)  NOT NULL COMMENT 'PAYMENT, REFUND',
    source_id       VARCHAR(64)  NOT NULL COMMENT 'payments.payment_id or refunds.refund_id from the payment service',
    payment_id      VARCHAR(64)  NOT NULL COMMENT 'Owning payments.payment_id from the payment service',
    amount          BIGINT       NOT NULL COMMENT 'Signed amount (smallest currency unit); refunds are negative',
    currency        CHAR(3)      NOT NULL COMMENT 'ISO 4217 currency code',
    payment_method  VARCHAR(32)  NULL     COMMENT 'Payment method (card, wallet, ...)',
    occurred_at     DATETIME(6)  NOT NULL COMMENT 'When the source payment/refund occurred',
    created_at      DATETIME(6)  NOT NULL,
    updated_at      DATETIME(6)  NULL,
    deleted_at      DATETIME(6)  NULL,

    UNIQUE KEY uk_settlement_detail_uid (uid),
    UNIQUE KEY uk_detail_source (source_type, source_id),
    INDEX idx_settlement_detail_date (settlement_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Individual payment/refund lines feeding a settlement bucket';


-- ------------------------------------------------------------
-- 3. Settlement job run log
--    Audit trail of settlement batch executions (Spring Batch).
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_job_run (
    id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid              BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    target_date      DATE          NOT NULL COMMENT 'Business date the batch run targeted',
    job_execution_id BIGINT        NULL     COMMENT 'Spring Batch BATCH_JOB_EXECUTION.JOB_EXECUTION_ID',
    status           VARCHAR(20)   NOT NULL COMMENT 'Run status (STARTED, COMPLETED, FAILED, ...)',
    started_at       DATETIME(6)   NOT NULL COMMENT 'Run start timestamp',
    finished_at      DATETIME(6)   NULL     COMMENT 'Run finish timestamp',
    message          VARCHAR(1024) NULL     COMMENT 'Human-readable outcome/error message',
    created_at       DATETIME(6)   NOT NULL,
    updated_at       DATETIME(6)   NULL,
    deleted_at       DATETIME(6)   NULL,

    UNIQUE KEY uk_settlement_job_run_uid (uid),
    INDEX idx_settlement_job_run_target (target_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Audit log of settlement batch executions';
