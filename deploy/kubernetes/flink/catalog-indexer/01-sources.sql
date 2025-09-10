CREATE TABLE publisher
(
    id         BIGINT NOT NULL,
    uid        STRING NOT NULL,
    name       STRING,
    address    STRING,
    ol_key     STRING,
    created_at TIMESTAMP(3),
    updated_at TIMESTAMP(3),
    deleted_at TIMESTAMP(3),
    PRIMARY KEY (id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.publisher',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE author
(
    id              BIGINT NOT NULL,
    uid             STRING NOT NULL,
    name            STRING,
    bio             STRING,
    birth_date      STRING,
    death_date      STRING,
    photo_uids      STRING,
    alternate_names STRING,
    ol_key          STRING,
    created_at      TIMESTAMP(3),
    updated_at      TIMESTAMP(3),
    deleted_at      TIMESTAMP(3),
    PRIMARY KEY (id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.author',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE work
(
    id                 BIGINT NOT NULL,
    uid                STRING NOT NULL,
    title              STRING,
    description        STRING,
    cover_uids         STRING,
    first_publish_date STRING,
    ol_key             STRING,
    created_at         TIMESTAMP(3),
    updated_at         TIMESTAMP(3),
    deleted_at         TIMESTAMP(3),
    PRIMARY KEY (id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.work',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE edition
(
    id              BIGINT NOT NULL,
    uid             STRING NOT NULL,
    title           STRING,
    isbn10          STRING,
    isbn13          STRING,
    number_of_pages INT,
    publish_date    STRING,
    cover_uids      STRING,
    languages       STRING,
    physical_format STRING,
    description     STRING,
    work_id         BIGINT NOT NULL,
    publisher_id    BIGINT,
    ol_key          STRING,
    created_at      TIMESTAMP(3),
    updated_at      TIMESTAMP(3),
    deleted_at      TIMESTAMP(3),
    PRIMARY KEY (id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.edition',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE subject
(
    id         BIGINT NOT NULL,
    uid        STRING NOT NULL,
    name       STRING,
    created_at TIMESTAMP(3),
    updated_at TIMESTAMP(3),
    deleted_at TIMESTAMP(3),
    PRIMARY KEY (id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.subject',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE work_author_mapping
(
    work_id   BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    PRIMARY KEY (work_id, author_id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.work_author_mapping',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );

CREATE TABLE work_subject_mapping
(
    work_id    BIGINT NOT NULL,
    subject_id BIGINT NOT NULL,
    PRIMARY KEY (work_id, subject_id) NOT ENFORCED
) WITH (
      'connector' = 'kafka',
      'topic' = 'catalog.public.work_subject_mapping',
      'properties.bootstrap.servers' = 'bookstore-kafka-kafka-bootstrap.kafka.svc.cluster.local:9093',
      'properties.group.id' = 'catalog-indexer',
      'scan.startup.mode' = 'earliest-offset',
      'value.format' = 'debezium-avro-confluent',
      'value.debezium-avro-confluent.url' = 'http://schemaregistry.kafka.svc.cluster.local:8081'
      );
