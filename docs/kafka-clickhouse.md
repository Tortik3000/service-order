# Kafka + ClickHouse pipeline

В `docker-compose.yml` добавлены сервисы `zookeeper`, `kafka`, `kafka-init` и `clickhouse`.

## Что настроено

1. В Kafka создаются топики:
   - `numbers` — основной поток с целыми числами.
   - `numbers-dlq` — DLQ для невалидных сообщений.
2. Сервис `kafka-init` отправляет стартовый набор сообщений в `numbers`:
   - `10`, `-4`, `3`, `-8`, `abc`, `15`.
3. В ClickHouse создан пайплайн:
   - `analytics.numbers_stream` (Kafka Engine) читает `numbers` как `RawBLOB`.
   - `analytics.mv_numbers_to_parsed` сохраняет валидные `Int64` в `analytics.numbers_parsed`.
   - `analytics.mv_numbers_to_dlq_topic` отправляет невалидные сообщения в Kafka топик `numbers-dlq`.
   - `analytics.numbers_dlq_stream` + `analytics.mv_dlq_topic_to_table` сохраняют DLQ сообщения в `analytics.numbers_dlq`.
   - `analytics.mv_sum_by_sign` считает суммы положительных и отрицательных чисел в `analytics.sum_by_sign`.

## Проверка

```bash
docker compose up -d zookeeper kafka kafka-init clickhouse
```

Проверить сумму:

```bash
docker exec -it service-order-clickhouse clickhouse-client -q "
SELECT
    sum(positive_sum) AS positive_sum,
    sum(negative_sum) AS negative_sum
FROM analytics.sum_by_sign;
"
```

Ожидаемо после стартовых сообщений:
- `positive_sum = 28` (`10 + 3 + 15`)
- `negative_sum = -12` (`-4 + -8`)

Проверить DLQ:

```bash
docker exec -it service-order-clickhouse clickhouse-client -q "
SELECT original_value, error, failed_at
FROM analytics.numbers_dlq
ORDER BY failed_at DESC
LIMIT 10;
"
```

Ожидаемо в DLQ будет запись для `abc`.
