package user

import (
	"fmt"

	"github.com/sing3demons/go-mongo-api/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	db *repository
}

func NewService(db *repository) *service {
	return &service{db}
}

func (s *service) Register(account User) (*User, error) {
	account.ID = primitive.NewObjectID()

	hashPassword, err := hashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	account.Password = string(hashPassword)
	// account.CreatedAt = primitive.NilObjectID.Timestamp()
	// account.UpdatedAt = primitive.NilObjectID.Timestamp()

	u, err := s.db.CreateAccount(account)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) FindAll() ([]User, error) {
	u, err := s.db.FindAll()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) FindOne(id string) (User, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error("can't convert id to object id err: ", err)
		return User{}, err
	}

	filter := bson.D{{Key: "_id", Value: _id}}
	fmt.Println(filter)
	user, err := s.db.FindOne(filter)
	if err != nil {
		return User{}, err
	}

	return *user, nil
}

func (s *service) Login(email, password string) (interface{}, error) {
	u, err := s.db.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(u.ID.Hex())
	if err != nil {
		logrus.Error("can't create token err: ", err)
		return nil, err
	}

	return token, nil
}

func hashPassword(password string) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashPassword, nil
}
