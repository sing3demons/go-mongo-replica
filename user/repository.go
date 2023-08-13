package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{db: db}
}

func (r *repository) collection() *mongo.Collection {
	return r.db.Collection("users")
}

func (r *repository) FindAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []User

	cursor, err := r.collection().Find(ctx, bson.M{})
	if err != nil {
		logrus.Error("can't find all user err: ", err)
		return nil, err
	}

	err = cursor.All(ctx, &users)
	if err != nil {
		logrus.Error("can't get all user err: ", err)
		return nil, err
	}

	return users, nil
}

func (r *repository) FindOne(fitter any) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User

	fmt.Println(fitter)

	err := r.collection().FindOne(ctx, fitter).Decode(&user)
	if err != nil {
		logrus.Error("can't find all user err: ", err)
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateAccount(account User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("+========")
	log.Println(account)

	_, err := r.collection().InsertOne(ctx, account)
	if err != nil {
		logrus.Error("can't create account err: ", err)
		return nil, err
	}

	return &account, nil
}
func (r *repository) FindByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User

	filter := bson.D{{Key: "email", Value: email}}
	err := r.collection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logrus.Error("can't find account err: ", err)
		return nil, err
	}

	return &user, nil
}
