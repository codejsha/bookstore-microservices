CREATE TABLE IF NOT EXISTS payments
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID        NOT NULL,
    payment_id      VARCHAR(64) NOT NULL,
    amount          BIGINT      NOT NULL,
    amount_captured BIGINT,
    currency        CHAR(3)     NOT NULL,
    status          VARCHAR(48) NOT NULL,
    payment_method  VARCHAR(32),
    captured_at     TIMESTAMP(6),
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    CONSTRAINT uk_payment_id UNIQUE (payment_id)
);

CREATE INDEX idx_payment_status ON payments (status);
CREATE INDEX idx_payment_captured ON payments (captured_at);

CREATE TABLE IF NOT EXISTS refunds
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID        NOT NULL,
    refund_id  VARCHAR(64) NOT NULL,
    payment_id VARCHAR(64) NOT NULL,
    amount     BIGINT      NOT NULL,
    currency   CHAR(3)     NOT NULL,
    status     VARCHAR(32) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    CONSTRAINT uk_refund_id UNIQUE (refund_id)
);

CREATE INDEX idx_refund_status ON refunds (status);
CREATE INDEX idx_refund_updated ON refunds (updated_at);
