SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS payments (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid               BINARY(16)  NOT NULL,
    payment_id        VARCHAR(64) NOT NULL,
    amount            BIGINT      NOT NULL,
    amount_captured   BIGINT      NULL,
    currency          CHAR(3)     NOT NULL,
    status            VARCHAR(48) NOT NULL,
    payment_method    VARCHAR(32) NULL,
    captured_at       DATETIME(6) NULL,
    created_at        DATETIME(6) NOT NULL,
    updated_at        DATETIME(6) NULL,
    deleted_at        DATETIME(6) NULL,

    UNIQUE KEY uk_payment_id (payment_id),
    INDEX idx_payment_status (status),
    INDEX idx_payment_captured (captured_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS refunds (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid               BINARY(16)  NOT NULL,
    refund_id         VARCHAR(64) NOT NULL,
    payment_id        VARCHAR(64) NOT NULL,
    amount            BIGINT      NOT NULL,
    currency          CHAR(3)     NOT NULL,
    status            VARCHAR(32) NOT NULL,
    created_at        DATETIME(6) NOT NULL,
    updated_at        DATETIME(6) NULL,
    deleted_at        DATETIME(6) NULL,

    UNIQUE KEY uk_refund_id (refund_id),
    INDEX idx_refund_status (status),
    INDEX idx_refund_updated (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
