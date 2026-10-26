-- ------------------------------------------------------------
-- Settlement run lock
--    Cluster-wide mutual exclusion for a settlement target date. One row per
--    business date; the holder is identified by owner_token (the run uid) and
--    keeps its lease alive through heartbeat_at/expires_at, so a crashed pod's
--    lock lapses on its own and the next run takes it over.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_run_lock
(
    id           BIGSERIAL PRIMARY KEY,
    uid          UUID         NOT NULL,
    target_date  DATE         NOT NULL,
    owner_token  VARCHAR(64)  NOT NULL,
    owner_host   VARCHAR(190) NOT NULL,
    acquired_at  TIMESTAMP(6) NOT NULL,
    heartbeat_at TIMESTAMP(6) NOT NULL,
    expires_at   TIMESTAMP(6) NOT NULL,
    released_at  TIMESTAMP(6),
    created_at   TIMESTAMP(6) NOT NULL,
    updated_at   TIMESTAMP(6),
    deleted_at   TIMESTAMP(6),
    CONSTRAINT uk_settlement_run_lock_uid UNIQUE (uid),
    CONSTRAINT uk_settlement_run_lock_date UNIQUE (target_date)
);

CREATE INDEX idx_settlement_run_lock_expires ON settlement_run_lock (expires_at);
