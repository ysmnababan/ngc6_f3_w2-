package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ngc5/model"
	"ngc5/pb"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	UserGRPC pb.UserServiceClient
}

func (h *UserHandler) AddUser(e echo.Context) error {
	var payload model.User

	err := e.Bind(&payload)
	if err != nil {
		log.Println("ERR BIND:", err)
		return e.JSON(http.StatusInternalServerError, "Error Binding")
	}

	//validate payload
	if payload.Name == "" {
		return e.JSON(http.StatusBadRequest, "Error or missing param")
	}

	in := &pb.AddRequest{
		Name: payload.Name,
	}

	response, err := h.UserGRPC.AddUser(context.TODO(), in)
	if err != nil {
		log.Println("ERR USER SERVICE: ", err)
		return e.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return e.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetUser(e echo.Context) error {
	in := &pb.GetRequest{
		Name: e.Get("name").(string),
	}

	res, err := h.UserGRPC.GetUser(context.TODO(), in)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// Extract the error code and handle accordingly
			switch st.Code() {
			case codes.NotFound:
				return e.JSON(http.StatusNotFound, st.Message())
			default:
				fmt.Printf("Unexpected error: %v\n", st.Message())
				return e.JSON(http.StatusInternalServerError, "Internal Server Error")
			}
		} else {
			fmt.Printf("Non-gRPC error: %v\n", err)
			log.Println("ERR USER SERVICE: ", err)
			return e.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
	}
	return e.JSON(http.StatusOK, res)
}

func generateToken(u model.User) (string, error) {
	// create the payload
	payload := jwt.MapClaims{
		"id":   u.ID,
		"name": u.Name,
		"exp":  time.Now().Add(time.Hour * 48).Unix(),
	}

	// define the method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// get token string
	_ = godotenv.Load()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("error when creating token: %v", err)
	}

	return tokenString, nil
}

func (h *UserHandler) Login(e echo.Context) error {
	var payload model.User

	err := e.Bind(&payload)
	if err != nil {
		log.Println("ERR BIND:", err)
		return e.JSON(http.StatusInternalServerError, "Error Binding")
	}

	//validate payload
	if payload.Name == "" {
		return e.JSON(http.StatusBadRequest, "Error or missing param")

	}

	in := &pb.GetRequest{
		Name: payload.Name,
	}
	res, err := h.UserGRPC.GetUser(context.TODO(), in)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// Extract the error code and handle accordingly
			switch st.Code() {
			case codes.NotFound:
				return e.JSON(http.StatusNotFound, st.Message())
			default:
				fmt.Printf("Unexpected error: %v\n", st.Message())
				return e.JSON(http.StatusInternalServerError, "Internal Server Error")
			}
		} else {
			fmt.Printf("Non-gRPC error: %v\n", err)
			log.Println("ERR USER SERVICE: ", err)
			return e.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	id, _ := primitive.ObjectIDFromHex(res.Id)
	user := model.User{
		Name: res.Name,
		ID:   id,
	}
	
	tokenString, err := generateToken(user)
	if err != nil {
		log.Println("ERR token: ", err)
		return e.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return e.JSON(http.StatusOK, tokenString)
}
