CREATE TABLE IF NOT EXISTS customer_ticket_category
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    parent_id   BIGINT,
    created_at  TIMESTAMP(6) NOT NULL,
    updated_at  TIMESTAMP(6),
    deleted_at  TIMESTAMP(6),
    actor       BIGINT       NOT NULL,
    version     BIGINT       NOT NULL,
    CONSTRAINT uk_ticket_category_uid UNIQUE (uid)
);

CREATE UNIQUE INDEX uq_ticket_category_name
    ON customer_ticket_category (name)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_ticket_category_parent_id ON customer_ticket_category (parent_id);

CREATE TABLE IF NOT EXISTS customer_ticket
(
    id           BIGSERIAL PRIMARY KEY,
    uid          UUID         NOT NULL,
    customer_uid UUID         NOT NULL,
    category_id  BIGINT,
    assignee_uid UUID,
    subject      VARCHAR(255) NOT NULL,
    description  TEXT         NOT NULL,
    status       VARCHAR(20)  NOT NULL,
    priority     VARCHAR(10)  NOT NULL,
    resolved_at  TIMESTAMP(6),
    created_at   TIMESTAMP(6) NOT NULL,
    updated_at   TIMESTAMP(6),
    deleted_at   TIMESTAMP(6),
    actor        BIGINT       NOT NULL,
    version      BIGINT       NOT NULL,
    CONSTRAINT uk_ticket_uid UNIQUE (uid)
);

CREATE INDEX idx_ticket_customer_uid ON customer_ticket (customer_uid);
CREATE INDEX idx_ticket_assignee_uid ON customer_ticket (assignee_uid);
CREATE INDEX idx_ticket_category_id ON customer_ticket (category_id);
CREATE INDEX idx_ticket_status ON customer_ticket (status);
CREATE INDEX idx_ticket_priority ON customer_ticket (priority);
CREATE INDEX idx_ticket_created_at ON customer_ticket (created_at);

CREATE TABLE IF NOT EXISTS customer_ticket_comment
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    ticket_id   BIGINT       NOT NULL,
    author_uid  UUID         NOT NULL,
    author_role VARCHAR(20)  NOT NULL,
    body        TEXT         NOT NULL,
    internal    BOOLEAN      NOT NULL,
    created_at  TIMESTAMP(6) NOT NULL,
    updated_at  TIMESTAMP(6),
    deleted_at  TIMESTAMP(6),
    actor       BIGINT       NOT NULL,
    version     BIGINT       NOT NULL,
    CONSTRAINT uk_ticket_comment_uid UNIQUE (uid)
);

CREATE INDEX idx_ticket_comment_ticket_id ON customer_ticket_comment (ticket_id);
CREATE INDEX idx_ticket_comment_created_at ON customer_ticket_comment (created_at);

CREATE TABLE IF NOT EXISTS customer_faq
(
    id          BIGSERIAL PRIMARY KEY,
    uid         UUID         NOT NULL,
    category_id BIGINT,
    question    VARCHAR(500) NOT NULL,
    answer      TEXT         NOT NULL,
    view_count  BIGINT       NOT NULL,
    published   BOOLEAN      NOT NULL,
    created_at  TIMESTAMP(6) NOT NULL,
    updated_at  TIMESTAMP(6),
    deleted_at  TIMESTAMP(6),
    actor       BIGINT       NOT NULL,
    version     BIGINT       NOT NULL,
    CONSTRAINT uk_faq_uid UNIQUE (uid)
);

CREATE INDEX idx_faq_category_id ON customer_faq (category_id);
CREATE INDEX idx_faq_published ON customer_faq (published);
