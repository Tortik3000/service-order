package user

import (
	"context"

	"github.com/Tortik3000/service-order/generated/api/user"
	"github.com/Tortik3000/service-order/internal/domain/entity"
)

type Handler interface {
	RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error)
}

type (
	userUseCase interface {
		RegisterUser(ctx context.Context, phone string, name string) (*entity.User, error)
	}
)

type handler struct {
	user.UnimplementedUserServiceServer
	uc userUseCase
}

var _ Handler = (*handler)(nil)

func NewUserHandler(u userUseCase) *handler {
	return &handler{uc: u}
}

func (h *handler) RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	u, err := h.uc.RegisterUser(ctx, req.Phone, req.Name)
	if err != nil {
		return nil, err
	}
	return &user.RegisterUserResponse{User: mapUserToProto(u)}, nil
}

func mapUserToProto(u *entity.User) *user.User {
	if u == nil {
		return nil
	}
	return &user.User{
		Id:    u.ID,
		Phone: u.Phone,
		Name:  u.Name,
	}
}
