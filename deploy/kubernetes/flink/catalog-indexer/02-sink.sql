CREATE TABLE book_search_index
(
    uid                STRING NOT NULL,
    title              STRING,
    description        STRING,
    ol_key             STRING,
    first_publish_date STRING,
    cover_uids         ARRAY<STRING>,
    created_at         TIMESTAMP_LTZ(3),
    updated_at         TIMESTAMP_LTZ(3),
    authors            ARRAY<ROW<
    uid             STRING,
    name               STRING,
    bio                STRING,
    birth_date         STRING,
    death_date         STRING,
    ol_key             STRING,
    alternate_names    ARRAY<STRING>
  >>,
    subjects           ARRAY<ROW<
    uid  STRING,
    name               STRING >>,
    editions           ARRAY<ROW<
    uid             STRING,
    title              STRING,
    isbn10             STRING,
    isbn13             STRING,
    publish_date       STRING,
    languages          ARRAY<STRING>,
    physical_format    STRING,
    description        STRING,
    publisher          ROW<
      uid  STRING,
    name               STRING >
  >>,
    PRIMARY KEY (uid) NOT ENFORCED
) WITH (
      'connector' = 'opensearch-2',
      'hosts' = 'http://opensearch-cluster-master.opensearch.svc.cluster.local:9200',
      'index' = 'books',
      'username' = '${OPENSEARCH_USERNAME}',
      'password' = '${OPENSEARCH_PASSWORD}',
      'sink.bulk-flush.max-actions' = '500',
      'sink.bulk-flush.interval' = '2s',
      'sink.delivery-guarantee' = 'at-least-once',
      'format' = 'json',
      'json.timestamp-format.standard' = 'ISO-8601'
      );
