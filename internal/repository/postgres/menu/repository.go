package menu

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tortik3000/service-order/pkg/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

const (
	categoryTable     = "menu_category"
	categoryID        = "id"
	categoryName      = "name"
	categorySortOrder = "sort_order"

	itemTable       = "menu_item"
	itemID          = "id"
	itemCategoryID  = "category_id"
	itemName        = "name"
	itemDescription = "description"
	itemPrice       = "price"
	itemActive      = "active"
	itemImageURL    = "image_url"
)

type Repository interface {
	GetCategory(ctx context.Context, id string) (*entity.Category, error)
	GetItemsByCategory(ctx context.Context, categoryID string) ([]entity.MenuItem, error)
	CreateCategory(ctx context.Context, category *entity.Category) error
	GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error)
	CreateMenuItem(ctx context.Context, item *entity.MenuItem) error
	UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error
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

func (r *repository) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	query := r.queryBuilder.
		Select(categoryID, categoryName, categorySortOrder).
		From(categoryTable).
		Where(sq.Eq{categoryID: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get category query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	cat := &entity.Category{}
	err = conn.QueryRow(ctx, sql, args...).Scan(&cat.ID, &cat.Name, &cat.SortOrder)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan category: %w", err)
	}

	return cat, nil
}

func (r *repository) GetItemsByCategory(ctx context.Context, catID string) ([]entity.MenuItem, error) {
	query := r.queryBuilder.
		Select(itemID, itemCategoryID, itemName, itemDescription, itemPrice, itemActive, itemImageURL).
		From(itemTable).
		Where(sq.Eq{itemCategoryID: catID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get items by category query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var items []entity.MenuItem
	for rows.Next() {
		var item entity.MenuItem
		err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.Active, &item.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) GetMenuItem(ctx context.Context, id string) (*entity.MenuItem, error) {
	query := r.queryBuilder.
		Select(itemID, itemCategoryID, itemName, itemDescription, itemPrice, itemActive, itemImageURL).
		From(itemTable).
		Where(sq.Eq{itemID: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get menu item query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	item := &entity.MenuItem{}
	err = conn.QueryRow(ctx, sql, args...).Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.Active, &item.ImageURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan menu item: %w", err)
	}

	return item, nil
}

func (r *repository) CreateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	query := r.queryBuilder.
		Insert(itemTable).
		Columns(itemCategoryID, itemName, itemDescription, itemPrice, itemActive, itemImageURL).
		Values(item.CategoryID, item.Name, item.Description, item.Price, item.Active, item.ImageURL).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build create menu item query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	err = conn.QueryRow(ctx, sql, args...).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("insert menu item: %w", err)
	}

	return nil
}

func (r *repository) UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	query := r.queryBuilder.
		Update(itemTable).
		Set(itemCategoryID, item.CategoryID).
		Set(itemName, item.Name).
		Set(itemDescription, item.Description).
		Set(itemPrice, item.Price).
		Set(itemActive, item.Active).
		Set(itemImageURL, item.ImageURL).
		Where(sq.Eq{itemID: item.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build update menu item query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update menu item: %w", err)
	}

	return nil
}

func (r *repository) CreateCategory(ctx context.Context, category *entity.Category) error {
	query := r.queryBuilder.
		Insert(categoryTable).
		Columns(categoryName, categorySortOrder).
		Values(category.Name, category.SortOrder).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build create category query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	err = conn.QueryRow(ctx, sql, args...).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}

	return nil
}
