CREATE TABLE IF NOT EXISTS customer_wishlist
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    user_id    BIGINT       NOT NULL,
    user_uid   UUID         NOT NULL,
    book_id    BIGINT       NOT NULL,
    book_uid   UUID         NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT       NOT NULL,
    version    BIGINT       NOT NULL,
    CONSTRAINT uk_wishlist_uid UNIQUE (uid),
    CONSTRAINT uq_customer_wishlist_user_book UNIQUE (user_id, book_id)
);

CREATE INDEX idx_wishlist_user_id ON customer_wishlist (user_id);

CREATE TABLE IF NOT EXISTS customer_point
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    user_id    BIGINT       NOT NULL,
    user_uid   UUID         NOT NULL,
    balance    INTEGER      NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT       NOT NULL,
    version    BIGINT       NOT NULL,
    CONSTRAINT uk_point_uid UNIQUE (uid),
    CONSTRAINT uq_customer_point_user UNIQUE (user_id)
);

CREATE TABLE IF NOT EXISTS customer_point_history
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    user_id     BIGINT       NOT NULL,
    user_uid    UUID         NOT NULL,
    change_type VARCHAR(20)  NOT NULL,
    amount      INTEGER      NOT NULL,
    reason      VARCHAR(255),
    created_at  TIMESTAMP(6) NOT NULL,
    CONSTRAINT chk_point_history_type
        CHECK (change_type IN ('EARN', 'SPEND', 'ADJUST', 'EXPIRE')),
    CONSTRAINT uk_point_history_uid UNIQUE (uid)
);

CREATE INDEX idx_point_history_user_id ON customer_point_history (user_id);

CREATE TABLE IF NOT EXISTS customer_review
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    user_id    BIGINT       NOT NULL,
    user_uid   UUID         NOT NULL,
    book_id    BIGINT       NOT NULL,
    book_uid   UUID         NOT NULL,
    rating     SMALLINT     NOT NULL,
    title      VARCHAR(255),
    content    TEXT,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT       NOT NULL,
    version    BIGINT       NOT NULL,
    CONSTRAINT chk_review_rating CHECK (rating BETWEEN 1 AND 5),
    CONSTRAINT uk_review_uid UNIQUE (uid),
    CONSTRAINT uq_review_user_book UNIQUE (user_id, book_id)
);

CREATE INDEX idx_review_user_id ON customer_review (user_id);
CREATE INDEX idx_review_book_id ON customer_review (book_id);
CREATE INDEX idx_review_rating ON customer_review (rating);
