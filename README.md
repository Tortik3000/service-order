# Площадка предварительного заказа онлайн в фастфуде

## 1. Техническое задание (лабораторная 1)

### 1.1 Назначение системы
Сервис предназначен для оформления предварительных заказов в фастфуде: клиент выбирает позицию меню, формирует заказ, указывает точку самовывоза и время получения. Система должна сократить время ожидания на точке выдачи.

### 1.2 Роли пользователей
1. **Клиент**
   - просматривает категории и позиции меню;
   - собирает заказ;
   - указывает номер телефона и точку самовывоза;
   - получает статус заказа.
2. **Оператор/администратор**
   - изменяет меню (категории и позиции);
   - отслеживает созданные заказы.

### 1.3 Функциональные требования
1. Управление клиентами:
   - создание клиента по номеру телефона;
   - получение клиента по ID.
2. Управление меню:
   - создание категории меню;
   - создание позиции меню (цена, активность, описание, категория);
   - чтение списка категорий/позиций.
3. Управление заказами:
   - создание заказа (клиент, список позиций, количество);
   - хранение статуса (`NEW`, `IN_PROGRESS`, `READY`, `CLOSED`, `CANCELED`);
   - расчет итоговой стоимости;
   - назначение точки выдачи и времени получения.
4. API и интеграции:
   - gRPC API;
   - HTTP-gateway для REST доступа;
   - Swagger/OpenAPI-спецификации.
5. Наблюдаемость:
   - метрики Prometheus;
   - дашборд Grafana;
   - алерты на недоступность сервиса и рост ошибок.
6. Аналитика потоков данных:
   - Kafka-топик `numbers` для входящих целых чисел;
   - DLQ-топик `numbers-dlq` для невалидных сообщений;
   - витрина в ClickHouse с суммой положительных/отрицательных значений.

### 1.4 Нефункциональные требования
- контейнеризация (Docker Compose);
- автоматическая миграция БД;
- CI/CD с публикацией Docker-образа и деплоем (GitHub Actions + Ansible);
- логирование и метрики для production-диагностики.

### 1.5 Ограничения
- MVP ориентирован на самовывоз (без доставки);
- хранение персональных данных ограничено номером телефона;
- базовая модель авторизации вне рамок лабораторной.

---

## 2. Прототип пользовательского интерфейса
Прототип описан в файле [docs/ui-prototype.md](docs/ui-prototype.md). Включены экраны:
- каталог меню;
- корзина и оформление;
- экран статуса заказа.

---

## 3. Схема базы данных и SQL-объекты

### 3.1 Логическая схема БД
Схема представлена в формате Mermaid в [docs/db-schema.md](docs/db-schema.md) (аналог диаграммы из Oracle SQL Developer Data Modeler).

### 3.2 Скрипты создания БД и объектов
- таблицы: `db/migrations/001..006_*.sql`;
- представление `order_summary_view` и хранимая функция `recalculate_order_total(...)`: `db/migrations/007_create_order_reporting_objects.sql`.

---

## 5. Реализация системы
Реализация выполнена на **Go** (допустимо условием «на любом языке»):
- запуск сервиса: `cmd/service-order/main.go`;
- use-case слой: `internal/usecase/*`;
- gRPC handlers: `internal/handlers/*`;
- PostgreSQL repositories: `internal/repository/postgres/*`;
- API контракты: `api/*.proto`, generated-код в `generated/`.

---

## 6. Развертывание и CI/CD

### Docker / docker-compose
- основной стек описан в `docker-compose.yml`;
- профили:
  - базовый (`service-order` + `postgres`),
  - `infra` (Kafka + ClickHouse),
  - `monitoring` (Prometheus + Grafana).

### CI/CD
- pipeline: `.github/workflows/deploy-with-ansible.yml`;
- шаги: сборка и публикация образа в Docker Hub, затем деплой через Ansible playbook.

---

## 10. Grafana + статус сервиса + уведомления

В проекте подготовлены:
- сбор метрик HTTP/gRPC и rate-limit (`pkg/metrics/*`);
- конфигурация Prometheus (`infra/prometheus/prometheus.yml`);
- provisioning Grafana datasource, dashboards и alert rules (`infra/grafana/provisioning/*`);
- готовый dashboard JSON (`infra/grafana/dashboards/dashboard.json`).

Настроенные уведомления:
1. `ServiceDown` — сервис недоступен;
2. `HighErrorRate` — всплеск ошибок на сервисе.

---

## 12. Интеграция Kafka и ClickHouse

Интеграция реализована в `infra/clickhouse/init/01_kafka_pipeline.sql`:
1. Kafka Engine-таблица `analytics.numbers_stream` читает топик `numbers`;
2. Materialized View `mv_numbers_to_parsed` парсит корректные `Int64` в `numbers_parsed`;
3. Materialized View `mv_numbers_to_dlq_topic` отправляет невалидные сообщения в DLQ (`numbers-dlq`);
4. `mv_dlq_topic_to_table` сохраняет DLQ в `numbers_dlq`;
5. `mv_sum_by_sign` обновляет агрегат `sum_by_sign` (суммы положительных и отрицательных чисел).

---

## Быстрый запуск

```bash
# базовый сервис
make docker-up

# c Kafka + ClickHouse
COMPOSE_PROFILES=infra docker compose up -d

# с мониторингом
COMPOSE_PROFILES=monitoring docker compose up -d
```
