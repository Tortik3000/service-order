package order

import (
	"context"

	"github.com/Tortik3000/service-order/generated/api/order"
	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Handler interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error)
	ListUserOrders(ctx context.Context, req *order.ListUserOrdersRequest) (*order.ListUserOrdersResponse, error)
	ListOrdersByStatus(ctx context.Context, req *order.ListOrdersByStatusRequest) (*order.ListOrdersByStatusResponse, error)
	UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error)
	CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error)
}

type (
	orderUseCase interface {
		CreateOrder(ctx context.Context, userID, restaurantID string, items []entity.OrderItem, pickUp bool) (*entity.Order, error)
		GetOrder(ctx context.Context, id string) (*entity.Order, error)
		ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]entity.Order, error)
		ListOrdersByStatus(ctx context.Context, statuses []entity.OrderStatus, limit, offset int32) ([]entity.Order, error)
		UpdateOrderStatus(ctx context.Context, id string, status entity.OrderStatus) (*entity.Order, error)
		CancelOrder(ctx context.Context, id string, reason string) (*entity.Order, error)
	}
)

type handler struct {
	order.UnimplementedOrderServiceServer
	uc orderUseCase
}

var _ Handler = (*handler)(nil)

func NewOrderHandler(u orderUseCase) *handler {
	return &handler{uc: u}
}

func (h *handler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	items := make([]entity.OrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.OrderItem{
			MenuItemID: it.MenuItemId,
			Quantity:   it.Quantity,
			UnitPrice:  it.UnitPrice,
		}
	}

	o, err := h.uc.CreateOrder(ctx, req.UserId, req.RestaurantId, items, req.PickUp)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Order: mapOrderToProto(o),
	}, nil
}

func (h *handler) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	o, err := h.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &order.GetOrderResponse{Order: mapOrderToProto(o)}, nil
}

func (h *handler) ListUserOrders(ctx context.Context, req *order.ListUserOrdersRequest) (*order.ListUserOrdersResponse, error) {
	orders, err := h.uc.ListUserOrders(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	res := make([]*order.Order, len(orders))
	for i, o := range orders {
		res[i] = mapOrderToProto(&o)
	}
	return &order.ListUserOrdersResponse{Orders: res}, nil
}

func (h *handler) ListOrdersByStatus(ctx context.Context, req *order.ListOrdersByStatusRequest) (*order.ListOrdersByStatusResponse, error) {
	statuses := make([]entity.OrderStatus, len(req.Statuses))
	for i, s := range req.Statuses {
		statuses[i] = entity.OrderStatus(s)
	}

	orders, err := h.uc.ListOrdersByStatus(ctx, statuses, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	res := make([]*order.Order, len(orders))
	for i, o := range orders {
		res[i] = mapOrderToProto(&o)
	}
	return &order.ListOrdersByStatusResponse{Orders: res}, nil
}

func (h *handler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	o, err := h.uc.UpdateOrderStatus(ctx, req.OrderId, entity.OrderStatus(req.NewStatus))
	if err != nil {
		return nil, err
	}
	return &order.UpdateOrderStatusResponse{Order: mapOrderToProto(o)}, nil
}

func (h *handler) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	o, err := h.uc.CancelOrder(ctx, req.OrderId, req.Reason)
	if err != nil {
		return nil, err
	}
	return &order.CancelOrderResponse{Order: mapOrderToProto(o)}, nil
}

func mapOrderToProto(o *entity.Order) *order.Order {
	items := make([]*order.OrderItem, len(o.Items))
	for i, it := range o.Items {
		items[i] = &order.OrderItem{
			MenuItemId: it.MenuItemID,
			Quantity:   it.Quantity,
			UnitPrice:  it.UnitPrice,
		}
	}

	return &order.Order{
		Id:           o.ID,
		UserId:       o.UserID,
		RestaurantId: o.RestaurantID,
		Status:       order.OrderStatus(o.Status),
		TotalAmount:  o.TotalAmount,
		Items:        items,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
		PickUp:       o.PickUp,
	}
}
