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
	CreatedAt    int64
}

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusDraft:
		return "draft"
	case OrderStatusAwaitingPayment:
		return "awaiting_payment"
	case OrderStatusPaid:
		return "paid"
	case OrderStatusInProgress:
		return "in_progress"
	case OrderStatusReady:
		return "ready"
	case OrderStatusCompleted:
		return "completed"
	case OrderStatusCancelled:
		return "cancelled"
	case OrderStatusFailed:
		return "failed"
	default:
		return "unspecified"
	}
}

func OrderStatusFromString(s string) OrderStatus {
	switch s {
	case "draft":
		return OrderStatusDraft
	case "awaiting_payment":
		return OrderStatusAwaitingPayment
	case "paid":
		return OrderStatusPaid
	case "in_progress":
		return OrderStatusInProgress
	case "ready":
		return OrderStatusReady
	case "completed":
		return OrderStatusCompleted
	case "cancelled":
		return OrderStatusCancelled
	case "failed":
		return OrderStatusFailed
	default:
		return OrderStatusUnspecified
	}
}
