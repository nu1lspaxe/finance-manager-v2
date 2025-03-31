package utils

import (
	"users/proto"

	"users/postgres/sqlc"
)

func SqlcToProto_FMUser(id int64, username string, email string, password string, createdAt, updatedAt int64) *proto.User {
	return &proto.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func SqlcToProto_FMAccount(account sqlc.FMAccount) *proto.UserAccount {
	return &proto.UserAccount{
		IdNumber: account.IDNumber,
		Balance:  account.Balance,
	}
}
