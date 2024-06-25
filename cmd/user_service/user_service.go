package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"ngc5/cmd/user_service/config"
	"ngc5/cmd/user_service/controller"
	"ngc5/cmd/user_service/repo"
	"ngc5/pb"

	"google.golang.org/grpc"
)

func main() {
	client, db := config.Connect("grpc")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	repo := &repo.Repo{DB: db}
	userController := &controller.UserController{UR: repo}

	// create new grpc server
	grpcServer := grpc.NewServer()

	// register the 'user' service server
	pb.RegisterUserServiceServer(grpcServer, userController)

	fmt.Println("USER MICROSERVICE")

	// start grpc server
	listen, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Println(err)
	}

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Println(err)
	}
	
}
