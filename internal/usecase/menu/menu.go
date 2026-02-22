package menu

import (
	"context"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Usecase interface {
	GetMenuByCategory(ctx context.Context, categoryID string) (*entity.Category, []entity.MenuItem, error)
	GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error)
	CreateMenuItem(ctx context.Context, item *entity.MenuItem) error
	UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error
	CreateCategory(ctx context.Context, category *entity.Category) error
}

type (
	menuRepository interface {
		GetCategory(ctx context.Context, id string) (*entity.Category, error)
		GetItemsByCategory(ctx context.Context, categoryID string) ([]entity.MenuItem, error)
		GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error)
		CreateMenuItem(ctx context.Context, item *entity.MenuItem) error
		UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error
		CreateCategory(ctx context.Context, category *entity.Category) error
	}
)

type useCase struct {
	menuRepo menuRepository
}

var _ Usecase = (*useCase)(nil)

func NewUseCase(menuRepo menuRepository) *useCase {
	return &useCase{menuRepo: menuRepo}
}

func (u *useCase) GetMenuByCategory(ctx context.Context, categoryID string) (*entity.Category, []entity.MenuItem, error) {
	cat, err := u.menuRepo.GetCategory(ctx, categoryID)
	if err != nil {
		return nil, nil, err
	}
	items, err := u.menuRepo.GetItemsByCategory(ctx, categoryID)
	if err != nil {
		return nil, nil, err
	}
	return cat, items, nil
}

func (u *useCase) GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error) {
	return u.menuRepo.GetMenuItem(ctx, id)
}

func (u *useCase) CreateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	return u.menuRepo.CreateMenuItem(ctx, item)
}

func (u *useCase) UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	return u.menuRepo.UpdateMenuItem(ctx, item)
}

func (u *useCase) CreateCategory(ctx context.Context, category *entity.Category) error {
	return u.menuRepo.CreateCategory(ctx, category)
}
