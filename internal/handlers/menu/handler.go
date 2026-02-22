package menu

import (
	"context"

	"github.com/Tortik3000/service-order/generated/api/menu"
	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Handler interface {
	GetMenuByCategory(ctx context.Context, req *menu.GetMenuByCategoryRequest) (*menu.GetMenuByCategoryResponse, error)
	GetMenuItem(ctx context.Context, req *menu.GetMenuItemRequest) (*menu.GetMenuItemResponse, error)
	CreateMenuItem(ctx context.Context, req *menu.CreateMenuItemRequest) (*menu.CreateMenuItemResponse, error)
	UpdateMenuItem(ctx context.Context, req *menu.UpdateMenuItemRequest) (*menu.UpdateMenuItemResponse, error)
	CreateCategory(ctx context.Context, req *menu.CreateCategoryRequest) (*menu.CreateCategoryResponse, error)
}

type (
	menuUseCase interface {
		GetMenuByCategory(ctx context.Context, categoryID string) (*entity.Category, []entity.MenuItem, error)
		GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error)
		CreateMenuItem(ctx context.Context, item *entity.MenuItem) error
		UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error
		CreateCategory(ctx context.Context, category *entity.Category) error
	}
)
type handler struct {
	menu.UnimplementedMenuServiceServer
	uc menuUseCase
}

var _ Handler = (*handler)(nil)

func NewMenuHandler(u menuUseCase) *handler {
	return &handler{uc: u}
}

func (h *handler) GetMenuByCategory(ctx context.Context, req *menu.GetMenuByCategoryRequest) (*menu.GetMenuByCategoryResponse, error) {
	cat, items, err := h.uc.GetMenuByCategory(ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}

	resItems := make([]*menu.MenuItem, len(items))
	for i, it := range items {
		resItems[i] = mapMenuItemToProto(&it)
	}

	return &menu.GetMenuByCategoryResponse{
		Category: mapCategoryToProto(cat),
		Items:    resItems,
	}, nil
}

func (h *handler) GetMenuItem(ctx context.Context, req *menu.GetMenuItemRequest) (*menu.GetMenuItemResponse, error) {
	it, err := h.uc.GetMenuItem(ctx, req.MenuItemId)
	if err != nil {
		return nil, err
	}
	return &menu.GetMenuItemResponse{Item: mapMenuItemToProto(it)}, nil
}

func (h *handler) CreateMenuItem(ctx context.Context, req *menu.CreateMenuItemRequest) (*menu.CreateMenuItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	it := &entity.MenuItem{
		CategoryID:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Active:      true,
	}
	if err := h.uc.CreateMenuItem(ctx, it); err != nil {
		return nil, err
	}
	return &menu.CreateMenuItemResponse{Item: mapMenuItemToProto(it)}, nil
}

func (h *handler) UpdateMenuItem(ctx context.Context, req *menu.UpdateMenuItemRequest) (*menu.UpdateMenuItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	it := &entity.MenuItem{
		ID:          req.Id,
		CategoryID:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Active:      true,
	}
	if err := h.uc.UpdateMenuItem(ctx, it); err != nil {
		return nil, err
	}
	return &menu.UpdateMenuItemResponse{Item: mapMenuItemToProto(it)}, nil
}

func (h *handler) CreateCategory(ctx context.Context, req *menu.CreateCategoryRequest) (*menu.CreateCategoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	cat := &entity.Category{
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}
	if err := h.uc.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return &menu.CreateCategoryResponse{Category: mapCategoryToProto(cat)}, nil
}

func mapCategoryToProto(c *entity.Category) *menu.Category {
	return &menu.Category{
		Id:        c.ID,
		Name:      c.Name,
		SortOrder: c.SortOrder,
	}
}

func mapMenuItemToProto(m *entity.MenuItem) *menu.MenuItem {
	return &menu.MenuItem{
		Id:          m.ID,
		CategoryId:  m.CategoryID,
		Name:        m.Name,
		Description: m.Description,
		Price:       int64(m.Price),
		Active:      m.Active,
	}
}
