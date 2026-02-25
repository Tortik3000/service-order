package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tortik3000/service-order/pkg/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

const (
	userTable = "customer"
	userID    = "id"
	userPhone = "phone"
	userName  = "name"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
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

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	query := r.queryBuilder.
		Insert(userTable).
		Columns(userPhone, userName).
		Values(user.Phone, user.Name).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build create user query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return err
	}

	err = conn.QueryRow(ctx, sql, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	query := r.queryBuilder.
		Select(userID, userPhone, userName).
		From(userTable).
		Where(sq.Eq{userPhone: phone})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get user by phone query: %w", err)
	}

	conn, err := r.transactor.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	err = conn.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Phone, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return user, nil
}
