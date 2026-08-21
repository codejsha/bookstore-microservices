CREATE TABLE IF NOT EXISTS payment
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id     BIGINT         NOT NULL,
    user_id      VARCHAR(36)    NOT NULL,
    payment_type VARCHAR(20)    NOT NULL,
    card_number  VARCHAR(20)    NOT NULL,
    amount       DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP      NOT NULL,
    created_at   TIMESTAMP      NOT NULL,
    updated_at   TIMESTAMP      NOT NULL,
    version      BIGINT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
