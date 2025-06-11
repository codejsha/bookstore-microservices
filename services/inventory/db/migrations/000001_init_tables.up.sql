CREATE TABLE IF NOT EXISTS warehouse
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    address    VARCHAR(255) NULL,
    capacity   INT NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    version    BIGINT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS stock
(
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    book_id           BIGINT NOT NULL,
    warehouse_id      BIGINT NOT NULL,
    quantity          INT    NOT NULL,
    reserved_quantity INT    NOT NULL,
    created_at        TIMESTAMP       NOT NULL,
    updated_at        TIMESTAMP       NOT NULL,
    version           BIGINT,
    FOREIGN KEY (warehouse_id) REFERENCES warehouse (id) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
