package order

import (
	"context"
	"fmt"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Usecase interface {
	CreateOrder(ctx context.Context, userID, restaurantID string, items []entity.OrderItem, pickUp bool) (*entity.Order, error)
	GetOrder(ctx context.Context, id string) (*entity.Order, error)
	ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error)
	ListOrdersByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status entity.OrderStatus) (*entity.Order, error)
	CancelOrder(ctx context.Context, id string, reason string) (*entity.Order, error)
}

type (
	orderRepository interface {
		Create(ctx context.Context, order *entity.Order) error
		CreateItems(ctx context.Context, orderID string, items []entity.OrderItem) error
		Get(ctx context.Context, id string) (*entity.Order, error)
		UpdateStatus(ctx context.Context, id string, status entity.OrderStatus) error
		ListByUser(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error)
		ListByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error)
	}

	menuRepository interface {
		GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error)
	}

	txManager interface {
		WithTx(ctx context.Context, function func(ctx context.Context) error) error
	}
)

type useCase struct {
	orderRepo  orderRepository
	menuRepo   menuRepository
	transactor txManager
}

var _ Usecase = (*useCase)(nil)

func NewUseCase(
	orderRepo orderRepository,
	menuRepo menuRepository,
	transactor txManager,
) *useCase {
	return &useCase{
		orderRepo:  orderRepo,
		menuRepo:   menuRepo,
		transactor: transactor,
	}
}

func (u *useCase) CreateOrder(ctx context.Context, userID, restaurantID string, items []entity.OrderItem, pickUp bool) (*entity.Order, error) {
	var totalAmount int64
	for i, item := range items {
		menuItem, err := u.menuRepo.GetMenuItem(ctx, item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("get menu item %s: %w", item.MenuItemID, err)
		}
		items[i].UnitPrice = int64(menuItem.Price * 100) // Сохраняем в копейках/центах
		totalAmount += items[i].UnitPrice * int64(item.Quantity)
	}

	order := &entity.Order{
		UserID:       userID,
		RestaurantID: restaurantID,
		Status:       entity.OrderStatusAwaitingPayment,
		TotalAmount:  totalAmount,
		Items:        items,
		PickUp:       pickUp,
	}

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := u.orderRepo.Create(ctx, order); err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		if err := u.orderRepo.CreateItems(ctx, order.ID, items); err != nil {
			return fmt.Errorf("create order items: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *useCase) GetOrder(ctx context.Context, id string) (*entity.Order, error) {
	return u.orderRepo.Get(ctx, id)
}

func (u *useCase) ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error) {
	return u.orderRepo.ListByUser(ctx, userID, limit, offset)
}

func (u *useCase) ListOrdersByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error) {
	return u.orderRepo.ListByStatus(ctx, statuses, limit, offset)
}

func (u *useCase) UpdateOrderStatus(ctx context.Context, id string, status entity.OrderStatus) (*entity.Order, error) {
	if err := u.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return u.orderRepo.Get(ctx, id)
}

func (u *useCase) CancelOrder(ctx context.Context, id string, reason string) (*entity.Order, error) {
	// В реальном приложении здесь могла бы быть логика отмены платежа
	return u.UpdateOrderStatus(ctx, id, entity.OrderStatusCancelled)
}
