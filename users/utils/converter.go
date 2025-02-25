package utils

import (
	"users/proto"

	"users/postgres/sqlc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func SqlcToProto(user sqlc.User) *proto.User {
	return &proto.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
	}
}
