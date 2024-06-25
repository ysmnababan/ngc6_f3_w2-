package main

import (
	"log"
	"ngc5/cmd/rest/handler"
	"ngc5/helper"
	"ngc5/pb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// create connection to grpc
	connection, err := grpc.NewClient(":50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	studentServiceClient := pb.NewUserServiceClient(connection)

	userHandler := &handler.UserHandler{UserGRPC: studentServiceClient}

	e := echo.New()

	e.Use(middleware.Recover())
	e.POST("/login", userHandler.Login)
	e.GET("/user", userHandler.GetUser, helper.Auth)
	e.POST("/user", userHandler.AddUser, helper.Auth)

	log.Fatal(e.Start(":8080"))
}
