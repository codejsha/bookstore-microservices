CREATE TABLE IF NOT EXISTS orders
(
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid             BINARY(16)     NOT NULL COMMENT 'Public UUID identifier',
    user_uid        BINARY(16)     NOT NULL COMMENT 'Owner user public UUID',
    order_number    VARCHAR(50)    NOT NULL,
    status          VARCHAR(20)    NOT NULL,
    currency        VARCHAR(3)     NOT NULL,
    items_amount    DECIMAL(19, 4) NOT NULL,
    discount_amount DECIMAL(19, 4) NOT NULL,
    shipping_amount DECIMAL(19, 4) NOT NULL,
    tax_amount      DECIMAL(19, 4) NOT NULL,
    total_amount    DECIMAL(19, 4) NOT NULL,
    idempotency_key VARCHAR(100)   NOT NULL,
    payment_uid     BINARY(16)     NULL COMMENT 'Public UUID of the captured payment',
    created_at      DATETIME(6)    NOT NULL,
    updated_at      DATETIME(6)    NULL,
    deleted_at      DATETIME(6)    NULL,
    actor_id        BIGINT         NOT NULL,
    version         BIGINT         NOT NULL,
    CONSTRAINT chk_orders_status
        CHECK (status IN ('PENDING', 'PAID', 'SHIPPED', 'DELIVERED', 'CANCELLED', 'REFUNDED')),
    CONSTRAINT chk_orders_currency CHECK (currency = UPPER(currency)),
    CONSTRAINT chk_orders_amounts
        CHECK (items_amount >= 0 AND discount_amount >= 0 AND shipping_amount >= 0 AND tax_amount >= 0 AND
               total_amount >= 0),
    UNIQUE KEY uk_uid (uid),
    UNIQUE KEY uq_orders_order_number (order_number),
    UNIQUE KEY uq_orders_idempotency_key (idempotency_key),
    INDEX idx_orders_user_uid (user_uid),
    INDEX idx_orders_status (status),
    INDEX idx_orders_created (created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS order_item
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid          BINARY(16)                                                   NOT NULL COMMENT 'Public UUID identifier',
    order_id     BIGINT                                                       NOT NULL,
    product_id   BIGINT                                                       NOT NULL,
    sku          VARCHAR(64)                                                  NULL,
    product_name VARCHAR(255)                                                 NULL,
    options      VARCHAR(255)                                                 NULL,
    quantity     INT                                                          NOT NULL,
    currency     VARCHAR(3)                                                   NOT NULL,
    price        DECIMAL(19, 4)                                               NOT NULL,
    tax_rate     DECIMAL(5, 4)                                                NOT NULL,
    subtotal     DECIMAL(19, 4) GENERATED ALWAYS AS (quantity * price) STORED NOT NULL,
    created_at   DATETIME(6)                                                  NOT NULL,
    updated_at   DATETIME(6)                                                  NULL,
    deleted_at   DATETIME(6)                                                  NULL,
    actor_id     BIGINT                                                       NOT NULL,
    version      BIGINT                                                       NOT NULL,
    CONSTRAINT fk_order_item__order
        FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
    CONSTRAINT chk_order_item_currency CHECK (currency = UPPER(currency)),
    CONSTRAINT chk_order_item_price CHECK (price >= 0),
    CONSTRAINT chk_order_item_subtotal CHECK (subtotal >= 0),
    UNIQUE KEY uk_order_item_uid (uid),
    INDEX idx_order_item_order_id (order_id),
    INDEX idx_order_item_product_id (product_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS order_adjustment
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid        BINARY(16)     NOT NULL COMMENT 'Public UUID identifier',
    order_id   BIGINT         NOT NULL,
    type       VARCHAR(30)    NOT NULL,
    label      VARCHAR(100)   NULL,
    amount     DECIMAL(19, 4) NOT NULL,
    meta       JSON           NULL,
    created_at DATETIME(6)    NOT NULL,
    updated_at DATETIME(6)    NULL,
    deleted_at DATETIME(6)    NULL,
    actor_id   BIGINT         NOT NULL,
    version    BIGINT         NOT NULL,
    CONSTRAINT fk_order_adjustment__order
        FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
    CONSTRAINT chk_order_adjustment_type CHECK (type IN ('COUPON', 'POINT', 'MANUAL', 'SHIPPING', 'TAX')),
    UNIQUE KEY uk_order_adjustment_uid (uid),
    INDEX idx_order_adjustment_order_id (order_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS order_shipping
(
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid             BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    order_id        BIGINT       NOT NULL,
    recipient_name  VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(20)  NOT NULL,
    address_line1   VARCHAR(255) NOT NULL,
    address_line2   VARCHAR(255) NULL,
    city            VARCHAR(100) NOT NULL,
    state           VARCHAR(100) NOT NULL,
    postal_code     VARCHAR(20)  NOT NULL,
    country         VARCHAR(2)   NOT NULL,
    shipping_method VARCHAR(100) NOT NULL,
    created_at      DATETIME(6)  NOT NULL,
    updated_at      DATETIME(6)  NULL,
    deleted_at      DATETIME(6)  NULL,
    actor_id        BIGINT       NOT NULL,
    version         BIGINT       NOT NULL,
    CONSTRAINT fk_order_shipping__order
        FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
    UNIQUE KEY uk_order_shipping_uid (uid),
    INDEX idx_order_shipping_order_id (order_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS cart
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid        BINARY(16)  NOT NULL COMMENT 'Public UUID identifier',
    user_uid   BINARY(16)  NOT NULL COMMENT 'Owner user public UUID',
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NULL,
    UNIQUE KEY uk_cart_uid (uid),
    UNIQUE KEY uq_cart_user (user_uid)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS cart_item
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid          BINARY(16)     NOT NULL COMMENT 'Public UUID identifier',
    cart_id      BIGINT         NOT NULL,
    product_id   BIGINT         NOT NULL,
    product_name VARCHAR(255)   NULL,
    quantity     INT            NOT NULL,
    currency     VARCHAR(3)     NOT NULL,
    price        DECIMAL(19, 4) NOT NULL,
    created_at   DATETIME(6)    NOT NULL,
    updated_at   DATETIME(6)    NULL,
    CONSTRAINT fk_cart_item__cart
        FOREIGN KEY (cart_id) REFERENCES cart (id) ON DELETE CASCADE,
    CONSTRAINT chk_cart_item_quantity CHECK (quantity > 0),
    CONSTRAINT chk_cart_item_price CHECK (price >= 0),
    UNIQUE KEY uk_cart_item_uid (uid),
    UNIQUE KEY uq_cart_item_product (cart_id, product_id),
    INDEX idx_cart_item_cart_id (cart_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
