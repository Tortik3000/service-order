# Схема базы данных (логическая)

```mermaid
erDiagram
    CUSTOMER {
        UUID id PK
        TEXT phone
    }

    PLACE {
        UUID id PK
        TEXT address
    }

    ORDERS {
        UUID id PK
        UUID customer_id FK
        TEXT status
        BIGINT total_amount
        TIME pickup_time
        UUID place_id FK
    }

    MENU_CATEGORY {
        UUID id PK
        TEXT name
        INT sort_order
    }

    MENU_ITEM {
        UUID id PK
        UUID category_id FK
        TEXT name
        TEXT description
        BIGINT price
        BOOLEAN active
    }

    ORDER_ITEM {
        UUID id PK
        UUID order_id FK
        UUID menu_item_id FK
        INT quantity
        BIGINT unit_price
    }

    CUSTOMER ||--o{ ORDERS : creates
    PLACE ||--o{ ORDERS : pickup_at
    MENU_CATEGORY ||--o{ MENU_ITEM : contains
    ORDERS ||--o{ ORDER_ITEM : has
    MENU_ITEM ||--o{ ORDER_ITEM : referenced
```

## Дополнительные объекты
- `order_summary_view` — агрегированная витрина по заказам (клиент, точка, сумма, количество позиций).
- `recalculate_order_total(order_id)` — пересчет суммы заказа по строкам `order_item`.
