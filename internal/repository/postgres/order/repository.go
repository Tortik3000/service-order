package order

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tortik3000/service-order/pkg/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

const (
	orderTable       = "orders"
	orderID          = "id"
	orderCustomerID  = "customer_id"
	orderStatus      = "status"
	orderTotalAmount = "total_amount"
	orderPickUp      = "pick_up"
	orderCreatedAt   = "created_at"
	orderUpdatedAt   = "updated_at"

	orderItemTable      = "order_item"
	orderItemOrderID    = "order_id"
	orderItemMenuItemID = "menu_item_id"
	orderItemQuantity   = "quantity"
	orderItemUnitPrice  = "unit_price"
)

type Repository interface {
	Create(ctx context.Context, order *entity.Order) error
	CreateItems(ctx context.Context, orderID string, items []entity.OrderItem) error
	Get(ctx context.Context, id string) (*entity.Order, error)
	UpdateStatus(ctx context.Context, id string, status entity.OrderStatus) error
	ListByUser(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error)
	ListByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error)
}

type (
	txManager interface {
		GetConn(ctx context.Context) (postgres.Conn, error)
	}
)

type repository struct {
	transactor   txManager
	queryBuilder sq.StatementBuilderType
}

var _ Repository = (*repository)(nil)

func New(transactor txManager) *repository {
	return &repository{
		transactor:   transactor,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *repository) Create(ctx context.Context, order *entity.Order) error {
	query := r.queryBuilder.
		Insert(orderTable).
		Columns(orderCustomerID, orderStatus, orderTotalAmount, orderPickUp).
		Values(order.UserID, order.Status, order.TotalAmount, order.PickUp).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s", orderID, orderCreatedAt, orderUpdatedAt))

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build create order query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	var createdAt, updatedAt time.Time
	err = conn.QueryRow(ctx, sql, args...).Scan(&order.ID, &createdAt, &updatedAt)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	order.CreatedAt = createdAt.Unix()
	order.UpdatedAt = updatedAt.Unix()

	return nil
}

func (r *repository) CreateItems(ctx context.Context, orderID string, items []entity.OrderItem) error {
	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	for _, item := range items {
		itemQuery := r.queryBuilder.
			Insert(orderItemTable).
			Columns(orderItemOrderID, orderItemMenuItemID, orderItemQuantity, orderItemUnitPrice).
			Values(orderID, item.MenuItemID, item.Quantity, item.UnitPrice)

		itemSql, itemArgs, err := itemQuery.ToSql()
		if err != nil {
			return fmt.Errorf("build create order item query: %w", err)
		}

		_, err = conn.Exec(ctx, itemSql, itemArgs...)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.Order, error) {
	query := r.queryBuilder.
		Select(orderID, orderCustomerID, orderStatus, orderTotalAmount, orderPickUp, orderCreatedAt, orderUpdatedAt).
		From(orderTable).
		Where(sq.Eq{orderID: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get order query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	order := &entity.Order{}
	var createdAt, updatedAt time.Time
	err = conn.QueryRow(ctx, sql, args...).Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount, &order.PickUp, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan order: %w", err)
	}
	order.CreatedAt = createdAt.Unix()
	order.UpdatedAt = updatedAt.Unix()

	itemsQuery := r.queryBuilder.
		Select(orderItemMenuItemID, orderItemQuantity, orderItemUnitPrice).
		From(orderItemTable).
		Where(sq.Eq{orderItemOrderID: id})

	itemsSql, itemsArgs, err := itemsQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get order items query: %w", err)
	}

	rows, err := conn.Query(ctx, itemsSql, itemsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.MenuItemID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status entity.OrderStatus) error {
	query := r.queryBuilder.
		Update(orderTable).
		Set(orderStatus, status).
		Set(orderUpdatedAt, "NOW()").
		Where(sq.Eq{orderID: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build update order status query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

func (r *repository) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error) {
	query := r.queryBuilder.
		Select(orderID, orderCustomerID, orderStatus, orderTotalAmount, orderPickUp, orderCreatedAt, orderUpdatedAt).
		From(orderTable).
		Where(sq.Eq{orderCustomerID: userID}).
		OrderBy(fmt.Sprintf("%s DESC", orderCreatedAt)).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list orders by user query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query orders by user: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount, &order.PickUp, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		order.CreatedAt = createdAt.Unix()
		order.UpdatedAt = updatedAt.Unix()
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) ListByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error) {
	if len(statuses) == 0 {
		return nil, nil
	}

	query := r.queryBuilder.
		Select(orderID, orderCustomerID, orderStatus, orderTotalAmount, orderPickUp, orderCreatedAt, orderUpdatedAt).
		From(orderTable).
		Where(sq.Eq{orderStatus: statuses}).
		OrderBy(fmt.Sprintf("%s DESC", orderCreatedAt)).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list orders by status query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query orders by status: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount, &order.PickUp, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		order.CreatedAt = createdAt.Unix()
		order.UpdatedAt = updatedAt.Unix()
		orders = append(orders, order)
	}

	return orders, nil
}
