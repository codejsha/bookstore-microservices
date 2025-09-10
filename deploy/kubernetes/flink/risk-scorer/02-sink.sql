CREATE TABLE risk_flag
(
    `key`   STRING,
    `value` STRING,
    PRIMARY KEY (`key`) NOT ENFORCED
) WITH (
      'connector' = 'redis',
      'host' = '${VALKEY_HOST}',
      'port' = '6379',
      'redis-mode' = 'single',
      'command' = 'set',
      'ttl' = '${AUTO_FLAG_TTL_SECONDS}'
      );
