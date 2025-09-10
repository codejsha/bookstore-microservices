SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS customers (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                 BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    customer_id         VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch customer_id (cus_xxx)',
    name                VARCHAR(255)  NULL     COMMENT 'Customer name',
    email               VARCHAR(255)  NULL     COMMENT 'Email address',
    phone               VARCHAR(32)   NULL     COMMENT 'Phone number',
    phone_country_code  VARCHAR(8)    NULL     COMMENT 'Phone country code (e.g. +82)',
    description         VARCHAR(500)  NULL     COMMENT 'Customer description',
    metadata            JSON          NULL     COMMENT 'Hyperswitch metadata (key-value)',
    default_billing_address  JSON     NULL     COMMENT 'Default billing address',
    default_shipping_address JSON     NULL     COMMENT 'Default shipping address',
    created_at          DATETIME(6)   NOT NULL,
    updated_at          DATETIME(6)   NULL,
    deleted_at          DATETIME(6)   NULL,

    UNIQUE KEY uk_uid (uid),
    UNIQUE KEY uk_customer_id (customer_id),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch customer information';


CREATE TABLE IF NOT EXISTS payment_methods (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                     BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    payment_method_id       VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch payment_method_id',
    customer_id             VARCHAR(64)   NOT NULL COMMENT 'FK -> customers.customer_id',
    payment_method          VARCHAR(32)   NOT NULL COMMENT 'card, wallet, bank_transfer, pay_later, etc.',
    payment_method_type     VARCHAR(32)   NULL     COMMENT 'credit, debit, google_pay, apple_pay, etc.',
    payment_method_issuer   VARCHAR(64)   NULL     COMMENT 'Issuer',
    card_network            VARCHAR(32)   NULL     COMMENT 'Visa, Mastercard, Amex, etc.',
    card_last4              CHAR(4)       NULL     COMMENT 'Last 4 digits of card number',
    card_exp_month          TINYINT UNSIGNED NULL  COMMENT 'Expiry month (1-12)',
    card_exp_year           SMALLINT UNSIGNED NULL COMMENT 'Expiry year (e.g. 2025)',
    card_holder_name        VARCHAR(255)  NULL     COMMENT 'Cardholder name',
    is_default              TINYINT(1)    NOT NULL COMMENT 'Whether this is the default payment method',
    metadata                JSON          NULL     COMMENT 'Additional metadata',
    created_at              DATETIME(6)   NOT NULL,
    updated_at              DATETIME(6)   NULL,
    deleted_at              DATETIME(6)   NULL,

    UNIQUE KEY uk_pm_uid (uid),
    UNIQUE KEY uk_payment_method_id (payment_method_id),
    INDEX idx_customer_id (customer_id),

    CONSTRAINT fk_pm_customer
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch saved payment methods';


-- ------------------------------------------------------------
-- 3. Mandates (Recurring payments)
--    Hyperswitch API: POST /payments (setup_future_usage)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS mandates (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                     BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    mandate_id              VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch mandate_id',
    customer_id             VARCHAR(64)   NOT NULL COMMENT 'FK -> customers.customer_id',
    payment_method_id       VARCHAR(64)   NULL     COMMENT 'FK -> payment_methods.payment_method_id',
    mandate_type            VARCHAR(32)   NOT NULL COMMENT 'single_use, multi_use',
    mandate_status          VARCHAR(32)   NOT NULL
                            COMMENT 'active, revoked, inactive, pending',
    -- Amount constraints
    mandate_amount          BIGINT        NULL     COMMENT 'Mandate amount (smallest currency unit, e.g. cents)',
    mandate_currency        CHAR(3)       NULL     COMMENT 'ISO 4217 currency code (e.g. USD, KRW)',
    start_date              DATETIME      NULL     COMMENT 'Mandate start date',
    end_date                DATETIME      NULL     COMMENT 'Mandate end date',
    -- CIT/MIT related
    setup_future_usage      VARCHAR(16)   NULL     COMMENT 'on_session, off_session',
    customer_acceptance_type VARCHAR(32)  NULL     COMMENT 'online, offline',
    customer_accepted_at    DATETIME      NULL     COMMENT 'Customer acceptance timestamp',
    metadata                JSON          NULL     COMMENT 'Additional metadata',
    created_at              DATETIME(6)   NOT NULL,
    updated_at              DATETIME(6)   NULL,
    deleted_at              DATETIME(6)   NULL,

    UNIQUE KEY uk_mandate_uid (uid),
    UNIQUE KEY uk_mandate_id (mandate_id),
    INDEX idx_mandate_customer (customer_id),
    INDEX idx_mandate_status (mandate_status),

    CONSTRAINT fk_mandate_customer
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_mandate_pm
        FOREIGN KEY (payment_method_id) REFERENCES payment_methods (payment_method_id)
        ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch mandate/recurring payment information';


-- ------------------------------------------------------------
-- 4. Payments
--    Hyperswitch API: POST /payments, /payments/{id}/confirm, /payments/{id}/capture
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payments (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                     BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    payment_id              VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch payment_id (pay_xxx)',
    idempotency_key         VARCHAR(100)  NULL     COMMENT 'Stable per-order-attempt key from the placement saga (order uid); NULL for REST-created payments',
    merchant_id             VARCHAR(64)   NULL     COMMENT 'Merchant identifier',
    profile_id              VARCHAR(64)   NULL     COMMENT 'Business profile ID',
    customer_id             VARCHAR(64)   NULL     COMMENT 'FK -> customers.customer_id',
    payment_method_id       VARCHAR(64)   NULL     COMMENT 'FK -> payment_methods.payment_method_id',
    mandate_id              VARCHAR(64)   NULL     COMMENT 'FK -> mandates.mandate_id (for recurring payments)',
    connector               VARCHAR(64)   NULL     COMMENT 'Payment processor (e.g. stripe, adyen)',
    connector_transaction_id VARCHAR(128) NULL     COMMENT 'Processor transaction ID',
    -- Amounts
    amount                  BIGINT        NOT NULL COMMENT 'Payment amount (smallest currency unit)',
    amount_capturable       BIGINT        NULL     COMMENT 'Capturable amount',
    amount_captured         BIGINT        NULL     COMMENT 'Captured amount',
    surcharge_amount        BIGINT        NULL     COMMENT 'Surcharge amount',
    tax_amount              BIGINT        NULL     COMMENT 'Tax amount',
    currency                CHAR(3)       NOT NULL COMMENT 'ISO 4217 currency code (e.g. USD, KRW)',
    -- Payment status & method
    status                  VARCHAR(48)   NOT NULL
                            COMMENT 'requires_payment_method, requires_confirmation, requires_customer_action, requires_capture, processing, succeeded, failed, cancelled, partially_captured, expired',
    capture_method          VARCHAR(20)   NOT NULL
                            COMMENT 'automatic, manual, manual_multiple, scheduled',
    authentication_type     VARCHAR(16)   NULL     COMMENT 'three_ds, no_three_ds',
    payment_method          VARCHAR(32)   NULL     COMMENT 'card, wallet, bank_transfer, etc.',
    payment_method_type     VARCHAR(32)   NULL     COMMENT 'credit, debit, google_pay, etc.',
    -- Session
    client_secret           VARCHAR(255)  NULL     COMMENT 'Client secret',
    setup_future_usage      VARCHAR(16)   NULL     COMMENT 'on_session, off_session',
    off_session             TINYINT(1)    NOT NULL COMMENT 'Merchant-initiated transaction flag',
    -- Additional information
    description             VARCHAR(500)  NULL     COMMENT 'Payment description',
    return_url              VARCHAR(1024) NULL     COMMENT 'Redirect URL',
    statement_descriptor    VARCHAR(255)  NULL     COMMENT 'Statement descriptor name',
    billing_address         JSON          NULL     COMMENT 'Billing address',
    shipping_address        JSON          NULL     COMMENT 'Shipping address',
    metadata                JSON          NULL     COMMENT 'Additional metadata',
    -- Errors
    error_code              VARCHAR(64)   NULL     COMMENT 'Unified error code',
    error_message           VARCHAR(1024) NULL     COMMENT 'Error message',
    -- Timestamps
    confirmed_at            DATETIME(6)   NULL     COMMENT 'Payment confirmation timestamp',
    captured_at             DATETIME(6)   NULL     COMMENT 'Capture timestamp',
    cancelled_at            DATETIME(6)   NULL     COMMENT 'Cancellation timestamp',
    created_at              DATETIME(6)   NOT NULL,
    updated_at              DATETIME(6)   NULL,
    deleted_at              DATETIME(6)   NULL,

    UNIQUE KEY uk_payment_uid (uid),
    UNIQUE KEY uk_payment_id (payment_id),
    UNIQUE KEY uq_payments_idempotency_key (idempotency_key),
    INDEX idx_payment_customer (customer_id),
    INDEX idx_payment_status (status),
    INDEX idx_payment_created (created_at),
    INDEX idx_payment_captured (captured_at),
    INDEX idx_payment_connector (connector, connector_transaction_id),
    INDEX idx_payment_mandate (mandate_id),

    CONSTRAINT fk_payment_customer
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_payment_pm
        FOREIGN KEY (payment_method_id) REFERENCES payment_methods (payment_method_id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_payment_mandate
        FOREIGN KEY (mandate_id) REFERENCES mandates (mandate_id)
        ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch payment information (Payment Intent)';


-- ------------------------------------------------------------
-- 5. Payment Attempts - Optional detail table
--    A single Payment can have multiple attempts (retries/routing)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payment_attempts (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                     BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    attempt_id              VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch attempt_id',
    payment_id              VARCHAR(64)   NOT NULL COMMENT 'FK -> payments.payment_id',
    connector               VARCHAR(64)   NULL     COMMENT 'Attempted processor',
    connector_transaction_id VARCHAR(128) NULL     COMMENT 'Processor transaction ID',
    amount                  BIGINT        NOT NULL COMMENT 'Attempt amount',
    currency                CHAR(3)       NOT NULL,
    status                  VARCHAR(48)   NOT NULL COMMENT 'Attempt status',
    authentication_type     VARCHAR(16)   NULL,
    payment_method          VARCHAR(32)   NULL,
    payment_method_type     VARCHAR(32)   NULL,
    error_code              VARCHAR(64)   NULL,
    error_message           VARCHAR(1024) NULL,
    created_at              DATETIME(6)   NOT NULL,
    updated_at              DATETIME(6)   NULL,
    deleted_at              DATETIME(6)   NULL,

    UNIQUE KEY uk_attempt_uid (uid),
    UNIQUE KEY uk_attempt_id (attempt_id),
    INDEX idx_attempt_payment (payment_id),

    CONSTRAINT fk_attempt_payment
        FOREIGN KEY (payment_id) REFERENCES payments (payment_id)
        ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Payment attempt history (smart routing/retries)';


-- ------------------------------------------------------------
-- 6. Refunds
--    Hyperswitch API: POST /refunds
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS refunds (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid                     BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    refund_id               VARCHAR(64)   NOT NULL COMMENT 'Hyperswitch refund_id',
    idempotency_key         VARCHAR(100)  NULL     COMMENT 'Stable per-payment key from the cancellation saga; NULL for REST-created refunds',
    payment_id              VARCHAR(64)   NOT NULL COMMENT 'FK -> payments.payment_id',
    connector               VARCHAR(64)   NULL     COMMENT 'Refund processor',
    connector_refund_id     VARCHAR(128)  NULL     COMMENT 'Processor refund ID',
    amount                  BIGINT        NOT NULL COMMENT 'Refund amount (smallest currency unit)',
    currency                CHAR(3)       NOT NULL COMMENT 'ISO 4217 currency code',
    status                  VARCHAR(32)   NOT NULL
                            COMMENT 'pending, succeeded, failed, review',
    reason                  VARCHAR(500)  NULL     COMMENT 'Refund reason',
    refund_type             VARCHAR(16)   NOT NULL
                            COMMENT 'instant, scheduled',
    error_code              VARCHAR(64)   NULL     COMMENT 'Unified error code',
    error_message           VARCHAR(1024) NULL     COMMENT 'Error message',
    metadata                JSON          NULL     COMMENT 'Additional metadata',
    created_at              DATETIME(6)   NOT NULL,
    updated_at              DATETIME(6)   NULL,
    deleted_at              DATETIME(6)   NULL,

    UNIQUE KEY uk_refund_uid (uid),
    UNIQUE KEY uk_refund_id (refund_id),
    UNIQUE KEY uq_refunds_idempotency_key (idempotency_key),
    INDEX idx_refund_payment (payment_id),
    INDEX idx_refund_status (status),
    INDEX idx_refund_created (created_at),

    CONSTRAINT fk_refund_payment
        FOREIGN KEY (payment_id) REFERENCES payments (payment_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch refund information';


-- ------------------------------------------------------------
-- 7. Webhook Event Log (Optional)
--    Stores incoming webhooks from Hyperswitch
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS webhook_events (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid               BINARY(16)    NOT NULL COMMENT 'Public UUID identifier',
    event_id          VARCHAR(64)   NULL     COMMENT 'Hyperswitch event_id',
    event_type        VARCHAR(64)   NOT NULL COMMENT 'payment_succeeded, refund_succeeded, etc.',
    object_type       VARCHAR(32)   NOT NULL COMMENT 'payment, refund, mandate, dispute',
    object_id         VARCHAR(64)   NOT NULL COMMENT 'Related object ID (e.g. payment_id)',
    payload           JSON          NOT NULL COMMENT 'Raw webhook JSON payload',
    signature         VARCHAR(512)  NULL     COMMENT 'HMAC signature',
    processed         TINYINT(1)    NOT NULL COMMENT 'Whether the event has been processed',
    created_at        DATETIME(6)   NOT NULL,
    updated_at        DATETIME(6)   NULL,
    deleted_at        DATETIME(6)   NULL,

    UNIQUE KEY uk_webhook_uid (uid),
    UNIQUE KEY uk_event_id (event_id),
    INDEX idx_webhook_object (object_type, object_id),
    INDEX idx_webhook_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Hyperswitch webhook event log';
