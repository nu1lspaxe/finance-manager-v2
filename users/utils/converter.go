package utils

import (
	"users/proto"

	"users/postgres/sqlc"
)

func SqlcToProto_FMUser(user sqlc.FMUser) *proto.User {
	return &proto.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	}
}

func SqlcToProto_FMAccount(account sqlc.FMAccount) *proto.UserAccount {
	return &proto.UserAccount{
		IdNumber: account.IDNumber,
		Balance:  account.Balance,
	}
}
