CREATE TABLE IF NOT EXISTS book_order
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id     VARCHAR(36)    NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status      VARCHAR(20)    NOT NULL,
    created_at  TIMESTAMP      NOT NULL,
    updated_at  TIMESTAMP      NOT NULL,
    version     BIGINT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS book_order_item
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id   BIGINT NOT NULL,
    book_id    BIGINT NOT NULL,
    quantity   INT    NOT NULL,
    created_at TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP       NOT NULL,
    version    BIGINT,
    FOREIGN KEY (order_id) REFERENCES book_order (id) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
