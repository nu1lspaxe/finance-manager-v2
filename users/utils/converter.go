package utils

import (
	"users/proto"

	"users/postgres/sqlc"
)

func SqlcToProto(user sqlc.FianaceManagerUser) *proto.User {
	return &proto.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	}
}
