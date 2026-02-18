CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.numbers_stream
(
    payload String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'numbers',
    kafka_group_name = 'clickhouse-numbers-consumer',
    kafka_format = 'RawBLOB',
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 0;

CREATE TABLE IF NOT EXISTS analytics.numbers_parsed
(
    value Int64,
    ingested_at DateTime DEFAULT now()
)
ENGINE = MergeTree
ORDER BY ingested_at;

CREATE TABLE IF NOT EXISTS analytics.sum_by_sign
(
    positive_sum Int64,
    negative_sum Int64
)
ENGINE = SummingMergeTree
ORDER BY tuple();

CREATE TABLE IF NOT EXISTS analytics.numbers_dlq_producer
(
    original_value String,
    error String,
    failed_at DateTime
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'numbers-dlq',
    kafka_format = 'JSONEachRow';

CREATE TABLE IF NOT EXISTS analytics.numbers_dlq_stream
(
    original_value String,
    error String,
    failed_at DateTime
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'numbers-dlq',
    kafka_group_name = 'clickhouse-dlq-consumer',
    kafka_format = 'JSONEachRow',
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 0;

CREATE TABLE IF NOT EXISTS analytics.numbers_dlq
(
    original_value String,
    error String,
    failed_at DateTime
)
ENGINE = MergeTree
ORDER BY failed_at;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.mv_numbers_to_parsed
TO analytics.numbers_parsed
AS
SELECT toInt64(payload) AS value, now() AS ingested_at
FROM analytics.numbers_stream
WHERE toInt64OrNull(payload) IS NOT NULL;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.mv_numbers_to_dlq_topic
TO analytics.numbers_dlq_producer
AS
SELECT
    payload AS original_value,
    'Message is not a valid Int64' AS error,
    now() AS failed_at
FROM analytics.numbers_stream
WHERE toInt64OrNull(payload) IS NULL;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.mv_dlq_topic_to_table
TO analytics.numbers_dlq
AS
SELECT
    original_value,
    error,
    failed_at
FROM analytics.numbers_dlq_stream;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.mv_sum_by_sign
TO analytics.sum_by_sign
AS
SELECT
    if(value > 0, value, 0) AS positive_sum,
    if(value < 0, value, 0) AS negative_sum
FROM analytics.numbers_parsed;
