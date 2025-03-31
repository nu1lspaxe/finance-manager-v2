package pkg

import (
	"context"
	"users/proto"
	"users/utils"

	"github.com/bufbuild/protovalidate-go"
)

type UserController struct {
	service *UserService
	proto.UnimplementedUserServiceServer
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	userId, token, err := c.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{UserId: userId, Token: token}, nil
}

func (c *UserController) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err = c.service.Logout(ctx, req.UserId, req.Token)
	if err != nil {
		return nil, err
	}
	return &proto.LogoutResponse{Success: true}, nil
}

func (c *UserController) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	user, err := c.service.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.UserResponse{User: utils.SqlcToProto_FMUser(user)}, nil
}

func (c *UserController) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	user, err := c.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.UserResponse{User: utils.SqlcToProto_FMUser(user)}, nil
}

func (c *UserController) GetAllUsers(_ *proto.EmptyRequest, stream proto.UserService_GetAllUsersServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), utils.TIMEOUT)
	defer cancel()

	userIds, err := c.service.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, userId := range userIds {
		err := stream.Send(&proto.UserIdResponse{UserId: userId})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *UserController) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err = c.service.UpdateUser(ctx, req.UserId, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateUserResponse{Success: true}, nil
}

func (c *UserController) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err := c.service.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteUserResponse{Success: true}, nil
}

func (c *UserController) AddUserAccount(ctx context.Context, req *proto.UserAccountRequest) (*proto.UserAccountResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	balance, err := c.service.AddAccount(ctx, req.UserId, req.IdNumber)
	if err != nil {
		return nil, err
	}
	return &proto.UserAccountResponse{Success: true, Balance: balance}, nil
}

func (c *UserController) GetUserAccounts(ctx context.Context, req *proto.GetUserAccountsRequest) (*proto.GetUserAccountsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	accounts, err := c.service.GetUserAccounts(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	protoAccounts := make([]*proto.UserAccount, 0, len(*accounts))
	for _, account := range *accounts {
		protoAccounts = append(protoAccounts, utils.SqlcToProto_FMAccount(account))
	}
	return &proto.GetUserAccountsResponse{Accounts: protoAccounts}, nil
}

func (c *UserController) DeleteAccount(ctx context.Context, req *proto.UserAccountRequest) (*proto.UserAccountResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err := c.service.DeleteAccount(ctx, req.IdNumber)
	if err != nil {
		return nil, err
	}
	return &proto.UserAccountResponse{Success: true}, nil
}
