CREATE TABLE IF NOT EXISTS publisher
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    name       VARCHAR(255) NOT NULL,
    address    VARCHAR(255),
    ol_key     VARCHAR(64),
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT       NOT NULL,
    version    BIGINT       NOT NULL,
    CONSTRAINT uk_publisher_uid UNIQUE (uid),
    CONSTRAINT uk_publisher_ol_key UNIQUE (ol_key)
);

CREATE INDEX idx_publisher_name ON publisher (name);

CREATE TABLE IF NOT EXISTS author
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    name            VARCHAR(255) NOT NULL,
    bio             TEXT,
    birth_date      VARCHAR(32),
    death_date      VARCHAR(32),
    photo_ids       JSONB,
    photo_uids      JSONB,
    alternate_names JSONB,
    ol_key          VARCHAR(64),
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    actor           BIGINT       NOT NULL,
    version         BIGINT       NOT NULL,
    CONSTRAINT uk_author_uid UNIQUE (uid),
    CONSTRAINT uk_author_ol_key UNIQUE (ol_key)
);

CREATE INDEX idx_author_name ON author (name);

CREATE TABLE IF NOT EXISTS work
(
    id                 BIGSERIAL PRIMARY KEY,
    uid                UUID         NOT NULL,
    title              VARCHAR(512) NOT NULL,
    description        TEXT,
    cover_ids          JSONB,
    cover_uids         JSONB,
    first_publish_date VARCHAR(32),
    ol_key             VARCHAR(64),
    created_at         TIMESTAMP(6) NOT NULL,
    updated_at         TIMESTAMP(6),
    deleted_at         TIMESTAMP(6),
    actor              BIGINT       NOT NULL,
    version            BIGINT       NOT NULL,
    CONSTRAINT uk_work_uid UNIQUE (uid),
    CONSTRAINT uk_work_ol_key UNIQUE (ol_key)
);

CREATE INDEX idx_work_title ON work (title);

CREATE TABLE IF NOT EXISTS edition
(
    id              BIGSERIAL PRIMARY KEY,
    uid             UUID         NOT NULL,
    title           VARCHAR(512) NOT NULL,
    isbn10          VARCHAR(10),
    isbn13          VARCHAR(13),
    number_of_pages INTEGER,
    publish_date    VARCHAR(32),
    cover_ids       JSONB,
    cover_uids      JSONB,
    languages       JSONB,
    physical_format VARCHAR(64),
    description     TEXT,
    work_id         BIGINT       NOT NULL,
    work_uid        UUID         NOT NULL,
    publisher_id    BIGINT,
    publisher_uid   UUID,
    ol_key          VARCHAR(64),
    created_at      TIMESTAMP(6) NOT NULL,
    updated_at      TIMESTAMP(6),
    deleted_at      TIMESTAMP(6),
    actor           BIGINT       NOT NULL,
    version         BIGINT       NOT NULL,
    CONSTRAINT fk_edition__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE RESTRICT,
    CONSTRAINT fk_edition__publisher
        FOREIGN KEY (publisher_id) REFERENCES publisher (id) ON DELETE SET NULL,
    CONSTRAINT chk_edition_isbn10_length
        CHECK (isbn10 IS NULL OR CHAR_LENGTH(isbn10) = 10),
    CONSTRAINT chk_edition_isbn13_length
        CHECK (isbn13 IS NULL OR CHAR_LENGTH(isbn13) = 13),
    CONSTRAINT uk_edition_uid UNIQUE (uid),
    CONSTRAINT uk_edition_ol_key UNIQUE (ol_key)
);

CREATE INDEX idx_edition_title ON edition (title);
CREATE INDEX idx_edition_isbn13 ON edition (isbn13);
CREATE INDEX idx_edition_work ON edition (work_id);
CREATE INDEX idx_edition_publisher ON edition (publisher_id);

CREATE TABLE IF NOT EXISTS subject
(
    id         BIGSERIAL PRIMARY KEY,
    uid        UUID         NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6),
    deleted_at TIMESTAMP(6),
    actor      BIGINT       NOT NULL,
    version    BIGINT       NOT NULL,
    CONSTRAINT uk_subject_uid UNIQUE (uid),
    CONSTRAINT uk_subject_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS work_author_mapping
(
    work_id   BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    PRIMARY KEY (work_id, author_id),
    CONSTRAINT fk_work_author__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE CASCADE,
    CONSTRAINT fk_work_author__author
        FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS work_subject_mapping
(
    work_id    BIGINT NOT NULL,
    subject_id BIGINT NOT NULL,
    PRIMARY KEY (work_id, subject_id),
    CONSTRAINT fk_work_subject__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE CASCADE,
    CONSTRAINT fk_work_subject__subject
        FOREIGN KEY (subject_id) REFERENCES subject (id) ON DELETE CASCADE
);
