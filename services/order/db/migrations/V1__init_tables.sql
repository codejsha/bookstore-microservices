CREATE TABLE IF NOT EXISTS orders
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID           NOT NULL,
    user_uid        UUID           NOT NULL,
    order_number    VARCHAR(50)    NOT NULL,
    status          VARCHAR(20)    NOT NULL,
    currency        VARCHAR(3)     NOT NULL,
    items_amount    DECIMAL(19, 4) NOT NULL,
    discount_amount DECIMAL(19, 4) NOT NULL,
    shipping_amount DECIMAL(19, 4) NOT NULL,
    tax_amount      DECIMAL(19, 4) NOT NULL,
    total_amount    DECIMAL(19, 4) NOT NULL,
    idempotency_key VARCHAR(100)   NOT NULL,
    payment_uid     UUID,
    created_at      TIMESTAMP(6)   NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    actor_id        BIGINT         NOT NULL,
    version         BIGINT         NOT NULL,
    CONSTRAINT uk_uid UNIQUE (uid),
    CONSTRAINT uq_orders_order_number UNIQUE (order_number),
    CONSTRAINT uq_orders_idempotency_key UNIQUE (idempotency_key)
);

CREATE INDEX idx_orders_user_uid ON orders (user_uid);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_orders_created ON orders (created_at);

CREATE TABLE IF NOT EXISTS order_item
(
    id           BIGSERIAL PRIMARY KEY,
    uid          UUID           NOT NULL,
    order_id     BIGINT         NOT NULL,
    product_id   BIGINT         NOT NULL,
    sku          VARCHAR(64),
    product_name VARCHAR(255),
    options      VARCHAR(255),
    quantity     INTEGER        NOT NULL,
    currency     VARCHAR(3)     NOT NULL,
    price        DECIMAL(19, 4) NOT NULL,
    tax_rate     DECIMAL(5, 4)  NOT NULL,
    subtotal     DECIMAL(19, 4) GENERATED ALWAYS AS (quantity * price) STORED NOT NULL,
    created_at   TIMESTAMP(6)   NOT NULL,
    updated_at   TIMESTAMP(6),
    deleted_at   TIMESTAMP(6),
    actor_id     BIGINT         NOT NULL,
    version      BIGINT         NOT NULL,
    CONSTRAINT uk_order_item_uid UNIQUE (uid)
);

CREATE INDEX idx_order_item_order_id ON order_item (order_id);
CREATE INDEX idx_order_item_product_id ON order_item (product_id);

CREATE TABLE IF NOT EXISTS order_adjustment
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID           NOT NULL,
    order_id   BIGINT         NOT NULL,
    type       VARCHAR(30)    NOT NULL,
    label      VARCHAR(100),
    amount     DECIMAL(19, 4) NOT NULL,
    meta       JSONB,
    created_at TIMESTAMP(6)   NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor_id   BIGINT         NOT NULL,
    version    BIGINT         NOT NULL,
    CONSTRAINT uk_order_adjustment_uid UNIQUE (uid)
);

CREATE INDEX idx_order_adjustment_order_id ON order_adjustment (order_id);

CREATE TABLE IF NOT EXISTS order_shipping
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    order_id        BIGINT       NOT NULL,
    recipient_name  VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(20)  NOT NULL,
    address_line1   VARCHAR(255) NOT NULL,
    address_line2   VARCHAR(255),
    city            VARCHAR(100) NOT NULL,
    state           VARCHAR(100) NOT NULL,
    postal_code     VARCHAR(20)  NOT NULL,
    country         VARCHAR(2)   NOT NULL,
    shipping_method VARCHAR(100) NOT NULL,
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    actor_id        BIGINT       NOT NULL,
    version         BIGINT       NOT NULL,
    CONSTRAINT uk_order_shipping_uid UNIQUE (uid)
);

CREATE INDEX idx_order_shipping_order_id ON order_shipping (order_id);

CREATE TABLE IF NOT EXISTS cart
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    user_uid   UUID         NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    CONSTRAINT uk_cart_uid UNIQUE (uid),
    CONSTRAINT uq_cart_user UNIQUE (user_uid)
);

CREATE TABLE IF NOT EXISTS cart_item
(
    id           BIGSERIAL PRIMARY KEY,
    uid          UUID           NOT NULL,
    cart_id      BIGINT         NOT NULL,
    product_id   BIGINT         NOT NULL,
    product_name VARCHAR(255),
    quantity     INTEGER        NOT NULL,
    currency     VARCHAR(3)     NOT NULL,
    price        DECIMAL(19, 4) NOT NULL,
    created_at   TIMESTAMP(6)   NOT NULL,
    updated_at   TIMESTAMP(6),
    CONSTRAINT uk_cart_item_uid UNIQUE (uid),
    CONSTRAINT uq_cart_item_product UNIQUE (cart_id, product_id)
);

CREATE INDEX idx_cart_item_cart_id ON cart_item (cart_id);
