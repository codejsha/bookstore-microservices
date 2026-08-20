CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    uid           UUID         NOT NULL,
    idp_uid       UUID         NOT NULL,
    email         VARCHAR(100) NOT NULL,
    first_name    VARCHAR(50)  NOT NULL,
    last_name     VARCHAR(50)  NOT NULL,
    phone         VARCHAR(30),
    roles         JSONB        NOT NULL,
    status        VARCHAR(20)  NOT NULL,
    last_login_at TIMESTAMP(6),
    created_at    TIMESTAMP(6) NOT NULL,
    updated_at    TIMESTAMP(6),
    deleted_at    TIMESTAMP(6),
    actor         BIGINT       NOT NULL,
    version       BIGINT       NOT NULL,
    CONSTRAINT chk_users_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DEACTIVATED')),
    CONSTRAINT uk_users_uid UNIQUE (uid),
    CONSTRAINT uk_users_idp_uid UNIQUE (idp_uid),
    CONSTRAINT uk_users_email UNIQUE (email)
);

CREATE INDEX idx_users_status ON users (status);
