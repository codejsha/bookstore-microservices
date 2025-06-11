CREATE TABLE IF NOT EXISTS wishlist
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    VARCHAR(36)     NOT NULL,
    book_id    BIGINT NOT NULL,
    created_at TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP       NOT NULL,
    version    BIGINT,
    UNIQUE (user_id, book_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS point
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    VARCHAR(36)  NOT NULL,
    point      INT NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    version    BIGINT,
    UNIQUE (user_id),
    INDEX idx_point_user_id (user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
