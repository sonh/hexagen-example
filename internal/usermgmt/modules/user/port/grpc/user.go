package grpc

import (
	"context"

	upb "backend/pkg/manabuf/usermgmt"
)

type UserService struct {
	upb.UnimplementedUserServiceServer
}

func (service *UserService) CreateUser(ctx context.Context, user *upb.User) (*upb.User, error) {
	// entity.ValidateUser(user)

	return nil, nil
}
