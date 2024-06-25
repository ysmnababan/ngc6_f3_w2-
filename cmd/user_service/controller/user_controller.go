package controller

import (
	"context"
	"fmt"
	"ngc5/cmd/user_service/repo"
	"ngc5/pb"
	"strings"

	"ngc5/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	UR repo.UserRepo
}

func (c *UserController) AddUser(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	user := model.User{
		Name: in.Name,
	}

	err := c.UR.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	resp := &pb.AddResponse{
		Name: user.Name,
		Id:   user.ID.Hex(),
	}
	return resp, nil
}

func (c *UserController) GetUser(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	user, err := c.UR.GetUser(in.Name)
	if err != nil {
		if strings.Contains(err.Error(), "no user in result set") {
			fmt.Println("here")
			return nil, status.Errorf(codes.NotFound, "no user found")
		}
		return nil, err
	}

	resp := &pb.GetResponse{
		Name: user.Name,
		Id:   user.ID.Hex(),
	}
	return resp, nil
}
