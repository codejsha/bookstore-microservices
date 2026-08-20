CREATE TABLE access_log
(
    `user_uid` STRING,
    `status`   INT,
    `path`     STRING,
    `method`   STRING,
    `kafka_ts` TIMESTAMP_LTZ(3) METADATA FROM 'timestamp',
    WATERMARK FOR `kafka_ts` AS `kafka_ts` - INTERVAL '30' SECOND
) WITH (
      'connector' = 'kafka',
      'topic' = 'bookstore.access-log',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'risk-scorer',
      'scan.startup.mode' = 'latest-offset',
      'format' = 'json',
      'json.ignore-parse-errors' = 'true'
      );
