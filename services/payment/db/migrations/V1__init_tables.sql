CREATE TABLE IF NOT EXISTS customers
(
    id                       BIGSERIAL PRIMARY KEY,
    uid                      UUID         NOT NULL,
    customer_id              VARCHAR(64)  NOT NULL,
    name                     VARCHAR(255),
    email                    VARCHAR(255),
    phone                    VARCHAR(32),
    phone_country_code       VARCHAR(8),
    description              VARCHAR(500),
    metadata                 JSONB,
    default_billing_address  JSONB,
    default_shipping_address JSONB,
    created_at               TIMESTAMP(6) NOT NULL,
    updated_at               TIMESTAMP(6),
    deleted_at               TIMESTAMP(6),
    CONSTRAINT uk_uid UNIQUE (uid),
    CONSTRAINT uk_customer_id UNIQUE (customer_id)
);

CREATE INDEX idx_email ON customers (email);

CREATE TABLE IF NOT EXISTS payment_methods
(
    id                    BIGSERIAL PRIMARY KEY,
    uid                   UUID         NOT NULL,
    payment_method_id     VARCHAR(64)  NOT NULL,
    customer_id           VARCHAR(64)  NOT NULL,
    payment_method        VARCHAR(32)  NOT NULL,
    payment_method_type   VARCHAR(32),
    payment_method_issuer VARCHAR(64),
    card_network          VARCHAR(32),
    card_last4            CHAR(4),
    card_exp_month        SMALLINT,
    card_exp_year         SMALLINT,
    card_holder_name      VARCHAR(255),
    is_default            BOOLEAN      NOT NULL,
    metadata              JSONB,
    created_at            TIMESTAMP(6) NOT NULL,
    updated_at            TIMESTAMP(6),
    deleted_at            TIMESTAMP(6),
    CONSTRAINT uk_pm_uid UNIQUE (uid),
    CONSTRAINT uk_payment_method_id UNIQUE (payment_method_id)
);

CREATE INDEX idx_customer_id ON payment_methods (customer_id);

CREATE TABLE IF NOT EXISTS mandates
(
    id                       BIGSERIAL PRIMARY KEY,
    uid                      UUID         NOT NULL,
    mandate_id               VARCHAR(64)  NOT NULL,
    customer_id              VARCHAR(64)  NOT NULL,
    payment_method_id        VARCHAR(64),
    mandate_type             VARCHAR(32)  NOT NULL,
    mandate_status           VARCHAR(32)  NOT NULL,
    mandate_amount           BIGINT,
    mandate_currency         CHAR(3),
    start_date               TIMESTAMP(6),
    end_date                 TIMESTAMP(6),
    setup_future_usage       VARCHAR(16),
    customer_acceptance_type VARCHAR(32),
    customer_accepted_at     TIMESTAMP(6),
    metadata                 JSONB,
    created_at               TIMESTAMP(6) NOT NULL,
    updated_at               TIMESTAMP(6),
    deleted_at               TIMESTAMP(6),
    CONSTRAINT uk_mandate_uid UNIQUE (uid),
    CONSTRAINT uk_mandate_id UNIQUE (mandate_id)
);

CREATE INDEX idx_mandate_customer ON mandates (customer_id);
CREATE INDEX idx_mandate_status ON mandates (mandate_status);
CREATE INDEX idx_mandate_pm ON mandates (payment_method_id);

CREATE TABLE IF NOT EXISTS payments
(
    id                       BIGSERIAL PRIMARY KEY,
    uid                      UUID          NOT NULL,
    payment_id               VARCHAR(64)   NOT NULL,
    idempotency_key          VARCHAR(100),
    merchant_id              VARCHAR(64),
    profile_id               VARCHAR(64),
    customer_id              VARCHAR(64),
    payment_method_id        VARCHAR(64),
    mandate_id               VARCHAR(64),
    connector                VARCHAR(64),
    connector_transaction_id VARCHAR(128),
    amount                   BIGINT        NOT NULL,
    amount_capturable        BIGINT,
    amount_captured          BIGINT,
    surcharge_amount         BIGINT,
    tax_amount               BIGINT,
    currency                 CHAR(3)       NOT NULL,
    status                   VARCHAR(48)   NOT NULL,
    capture_method           VARCHAR(20)   NOT NULL,
    authentication_type      VARCHAR(16),
    payment_method           VARCHAR(32),
    payment_method_type      VARCHAR(32),
    client_secret            VARCHAR(255),
    setup_future_usage       VARCHAR(16),
    off_session              BOOLEAN       NOT NULL,
    description              VARCHAR(500),
    return_url               VARCHAR(1024),
    statement_descriptor     VARCHAR(255),
    billing_address          JSONB,
    shipping_address         JSONB,
    metadata                 JSONB,
    error_code               VARCHAR(64),
    error_message            VARCHAR(1024),
    confirmed_at             TIMESTAMP(6),
    captured_at              TIMESTAMP(6),
    cancelled_at             TIMESTAMP(6),
    created_at               TIMESTAMP(6)  NOT NULL,
    updated_at               TIMESTAMP(6),
    deleted_at               TIMESTAMP(6),
    CONSTRAINT uk_payment_uid UNIQUE (uid),
    CONSTRAINT uk_payment_id UNIQUE (payment_id),
    CONSTRAINT uq_payments_idempotency_key UNIQUE (idempotency_key)
);

CREATE INDEX idx_payment_customer ON payments (customer_id);
CREATE INDEX idx_payment_status ON payments (status);
CREATE INDEX idx_payment_created ON payments (created_at);
CREATE INDEX idx_payment_captured ON payments (captured_at);
CREATE INDEX idx_payment_connector ON payments (connector, connector_transaction_id);
CREATE INDEX idx_payment_mandate ON payments (mandate_id);
CREATE INDEX idx_payment_pm ON payments (payment_method_id);

CREATE TABLE IF NOT EXISTS payment_attempts
(
    id                       BIGSERIAL PRIMARY KEY,
    uid                      UUID          NOT NULL,
    attempt_id               VARCHAR(64)   NOT NULL,
    payment_id               VARCHAR(64)   NOT NULL,
    connector                VARCHAR(64),
    connector_transaction_id VARCHAR(128),
    amount                   BIGINT        NOT NULL,
    currency                 CHAR(3)       NOT NULL,
    status                   VARCHAR(48)   NOT NULL,
    authentication_type      VARCHAR(16),
    payment_method           VARCHAR(32),
    payment_method_type      VARCHAR(32),
    error_code               VARCHAR(64),
    error_message            VARCHAR(1024),
    created_at               TIMESTAMP(6)  NOT NULL,
    updated_at               TIMESTAMP(6),
    deleted_at               TIMESTAMP(6),
    CONSTRAINT uk_attempt_uid UNIQUE (uid),
    CONSTRAINT uk_attempt_id UNIQUE (attempt_id)
);

CREATE INDEX idx_attempt_payment ON payment_attempts (payment_id);

CREATE TABLE IF NOT EXISTS refunds
(
    id                  BIGSERIAL PRIMARY KEY,
    uid                 UUID          NOT NULL,
    refund_id           VARCHAR(64)   NOT NULL,
    idempotency_key     VARCHAR(100),
    payment_id          VARCHAR(64)   NOT NULL,
    connector           VARCHAR(64),
    connector_refund_id VARCHAR(128),
    amount              BIGINT        NOT NULL,
    currency            CHAR(3)       NOT NULL,
    status              VARCHAR(32)   NOT NULL,
    reason              VARCHAR(500),
    refund_type         VARCHAR(16)   NOT NULL,
    error_code          VARCHAR(64),
    error_message       VARCHAR(1024),
    metadata            JSONB,
    created_at          TIMESTAMP(6)  NOT NULL,
    updated_at          TIMESTAMP(6),
    deleted_at          TIMESTAMP(6),
    CONSTRAINT uk_refund_uid UNIQUE (uid),
    CONSTRAINT uk_refund_id UNIQUE (refund_id),
    CONSTRAINT uq_refunds_idempotency_key UNIQUE (idempotency_key)
);

CREATE INDEX idx_refund_payment ON refunds (payment_id);
CREATE INDEX idx_refund_status ON refunds (status);
CREATE INDEX idx_refund_created ON refunds (created_at);

CREATE TABLE IF NOT EXISTS webhook_events
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    event_id    VARCHAR(64),
    event_type  VARCHAR(64)  NOT NULL,
    object_type VARCHAR(32)  NOT NULL,
    object_id   VARCHAR(64)  NOT NULL,
    payload     JSONB        NOT NULL,
    signature   VARCHAR(512),
    processed   BOOLEAN      NOT NULL,
    created_at  TIMESTAMP(6) NOT NULL,
    updated_at  TIMESTAMP(6),
    deleted_at  TIMESTAMP(6),
    CONSTRAINT uk_webhook_uid UNIQUE (uid),
    CONSTRAINT uk_event_id UNIQUE (event_id)
);

CREATE INDEX idx_webhook_object ON webhook_events (object_type, object_id);
CREATE INDEX idx_webhook_processed ON webhook_events (processed, created_at);
