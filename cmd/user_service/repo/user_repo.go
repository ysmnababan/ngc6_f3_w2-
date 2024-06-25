package repo

import (
	"context"
	"errors"
	"log"
	"ngc5/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	GetUser(name string) (model.User, error)
	CreateUser(user *model.User) error
}

type Repo struct {
	DB *mongo.Database
}

func (r *Repo) GetUser(name string) (model.User, error) {
	var user model.User
	err := r.DB.Collection("Users").FindOne(context.TODO(), bson.M{"name": name}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.User{}, errors.New("no user in result set")
		}
		log.Println("ERR:", err)
		return model.User{}, err
	}
	return user, nil
}

func (r *Repo) CreateUser(user *model.User) error {
	res, err := r.DB.Collection("Users").InsertOne(context.TODO(), *user)
	if err != nil {
		log.Println("ERR:", err)
		return err
	}

	// Check if the InsertedID is of type primitive.ObjectID
	if objectID, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = objectID
	}
	return nil
}
