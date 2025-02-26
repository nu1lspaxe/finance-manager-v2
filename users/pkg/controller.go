package pkg

import (
	"context"
	"users/proto"
	"users/utils"
)

type UserController struct {
	service *UserService
	proto.UnimplementedUserServiceServer
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) CreateUser(context context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user, err := c.service.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.CreateUserResponse{User: utils.SqlcToProto(*user)}, nil
}

func (c *UserController) GetUser(context context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := c.service.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{User: utils.SqlcToProto(*user)}, nil
}

func (c *UserController) ListUsers(context context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	users, err := c.service.ListUsers(req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*proto.User, 0, len(*users))
	for _, user := range *users {
		protoUsers = append(protoUsers, utils.SqlcToProto(user))
	}
	return &proto.ListUsersResponse{Users: protoUsers}, nil
}

func (c *UserController) UpdateUser(context context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	err := c.service.UpdateUser(req.Id, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateUserResponse{Success: true}, nil
}

func (c *UserController) DeleteUser(context context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := c.service.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteUserResponse{Success: true}, nil
}
