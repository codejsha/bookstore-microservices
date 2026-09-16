SET NAMES utf8mb4;

-- ------------------------------------------------------------
-- Settlement run lock
--    Cluster-wide mutual exclusion for a settlement target date. One row per
--    business date; the holder is identified by owner_token (the run uid) and
--    keeps its lease alive through heartbeat_at/expires_at, so a crashed pod's
--    lock lapses on its own and the next run takes it over.
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS settlement_run_lock (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uid          BINARY(16)   NOT NULL COMMENT 'Public UUID identifier',
    target_date  DATE         NOT NULL COMMENT 'Business date the lock guards',
    owner_token  VARCHAR(64)  NOT NULL COMMENT 'Run uid of the current holder',
    owner_host   VARCHAR(190) NOT NULL COMMENT 'Pod/host name of the current holder',
    acquired_at  DATETIME(6)  NOT NULL COMMENT 'When the current holder took the lock',
    heartbeat_at DATETIME(6)  NOT NULL COMMENT 'Last liveness renewal by the holder',
    expires_at   DATETIME(6)  NOT NULL COMMENT 'Lease deadline; a lapsed lease may be taken over',
    released_at  DATETIME(6)  NULL     COMMENT 'When the holder released the lock; NULL while held',
    created_at   DATETIME(6)  NOT NULL,
    updated_at   DATETIME(6)  NULL,
    deleted_at   DATETIME(6)  NULL,

    UNIQUE KEY uk_settlement_run_lock_uid (uid),
    UNIQUE KEY uk_settlement_run_lock_date (target_date),
    INDEX idx_settlement_run_lock_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Per-target-date mutual exclusion for settlement runs';
