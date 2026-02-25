package entity

type OrderStatus int32

const (
	OrderStatusUnspecified OrderStatus = iota
	OrderStatusDraft
	OrderStatusAwaitingPayment
	OrderStatusPaid
	OrderStatusInProgress
	OrderStatusReady
	OrderStatusCompleted
	OrderStatusCancelled
	OrderStatusFailed
)

type OrderItem struct {
	MenuItemID string
	Quantity   int32
	UnitPrice  int64
}

type Order struct {
	ID           string
	UserID       string
	RestaurantID string
	Status       OrderStatus
	TotalAmount  int64
	Items        []OrderItem
	PickUp       bool
	CreatedAt    int64
	UpdatedAt    int64
}
