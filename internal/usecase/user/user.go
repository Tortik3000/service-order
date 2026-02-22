package user

import (
	"context"

	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Usecase interface {
	RegisterUser(ctx context.Context, phone string) (*entity.User, error)
}

type (
	userRepository interface {
		Create(ctx context.Context, user *entity.User) error
		GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	}
)

type useCase struct {
	userRepo userRepository
}

var _ Usecase = (*useCase)(nil)

func NewUseCase(userRepo userRepository) *useCase {
	return &useCase{userRepo: userRepo}
}

func (u *useCase) RegisterUser(ctx context.Context, phone string) (*entity.User, error) {
	user, err := u.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user = &entity.User{
		Phone: phone,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
