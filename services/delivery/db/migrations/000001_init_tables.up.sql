CREATE TABLE IF NOT EXISTS carrier
(
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid           BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    name          VARCHAR(255) NOT NULL,
    code          VARCHAR(50)  NOT NULL,
    contact_name  VARCHAR(255) NULL,
    contact_phone VARCHAR(50)  NULL,
    contact_email VARCHAR(255) NULL,
    base_rate     DOUBLE       NOT NULL,
    rate_per_kg   DOUBLE       NOT NULL,
    status        VARCHAR(10)  NOT NULL COMMENT 'ACTIVE, INACTIVE',
    created_at    DATETIME(6)  NOT NULL,
    updated_at    DATETIME(6)  NOT NULL,
    deleted_at    DATETIME(6)  NULL,
    CONSTRAINT chk_carrier_status CHECK (status IN ('ACTIVE', 'INACTIVE')),
    UNIQUE KEY uk_carrier_uid (uid),
    UNIQUE KEY uk_carrier_code (code),
    INDEX idx_carrier_name (name),
    INDEX idx_carrier_status (status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shipment
(
    id                       BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid                      BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    order_uid                BINARY(16)   NOT NULL,
    carrier_uid              BINARY(16)   NULL,
    origin_address           VARCHAR(500) NOT NULL,
    destination_address      VARCHAR(500) NOT NULL,
    destination_city         VARCHAR(100) NOT NULL,
    destination_state        VARCHAR(100) NOT NULL,
    destination_country_code VARCHAR(3)   NOT NULL,
    destination_postal_code  VARCHAR(20)  NOT NULL,
    status                   VARCHAR(20)  NOT NULL COMMENT 'PLANNED, DISPATCHED, PICKED_UP, IN_TRANSIT, OUT_FOR_DELIVERY, DELIVERED, FAILED, CANCELLED',
    tracking_number          VARCHAR(100) NULL,
    weight_kg                DOUBLE       NULL,
    planned_pickup_at        DATETIME(6)  NULL,
    planned_delivery_at      DATETIME(6)  NULL,
    actual_pickup_at         DATETIME(6)  NULL,
    actual_delivery_at       DATETIME(6)  NULL,
    created_at               DATETIME(6)  NOT NULL,
    updated_at               DATETIME(6)  NOT NULL,
    deleted_at               DATETIME(6)  NULL,
    CONSTRAINT chk_shipment_status
        CHECK (status IN ('PLANNED', 'DISPATCHED', 'PICKED_UP', 'IN_TRANSIT', 'OUT_FOR_DELIVERY', 'DELIVERED', 'FAILED', 'CANCELLED')),
    UNIQUE KEY uk_shipment_uid (uid),
    INDEX idx_shipment_order_uid (order_uid),
    INDEX idx_shipment_carrier_uid (carrier_uid),
    INDEX idx_shipment_status (status),
    INDEX idx_shipment_created (created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tracking
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid          BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    shipment_uid BINARY(16)   NOT NULL,
    status       VARCHAR(20)  NOT NULL,
    location     VARCHAR(255) NOT NULL,
    description  TEXT         NOT NULL,
    occurred_at  DATETIME(6)  NOT NULL,
    created_at   DATETIME(6)  NOT NULL,
    UNIQUE KEY uk_tracking_uid (uid),
    INDEX idx_tracking_shipment_uid (shipment_uid),
    INDEX idx_tracking_occurred (occurred_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS freight
(
    id                 BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid                BINARY(16)  NOT NULL COMMENT 'Public UUID identifier',
    shipment_uid       BINARY(16)  NOT NULL,
    carrier_uid        BINARY(16)  NOT NULL,
    base_cost          DOUBLE      NOT NULL,
    weight_surcharge   DOUBLE      NOT NULL DEFAULT 0,
    distance_surcharge DOUBLE      NOT NULL DEFAULT 0,
    discount           DOUBLE      NOT NULL DEFAULT 0,
    total_cost         DOUBLE      NOT NULL,
    currency           VARCHAR(3)  NOT NULL DEFAULT 'KRW',
    status             VARCHAR(10) NOT NULL COMMENT 'ESTIMATED, CONFIRMED, INVOICED, PAID',
    invoiced_at        DATETIME(6) NULL,
    paid_at            DATETIME(6) NULL,
    created_at         DATETIME(6) NOT NULL,
    updated_at         DATETIME(6) NOT NULL,
    deleted_at         DATETIME(6) NULL,
    CONSTRAINT chk_freight_status CHECK (status IN ('ESTIMATED', 'CONFIRMED', 'INVOICED', 'PAID')),
    UNIQUE KEY uk_freight_uid (uid),
    INDEX idx_freight_shipment_uid (shipment_uid),
    INDEX idx_freight_carrier_uid (carrier_uid),
    INDEX idx_freight_status (status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
