package user

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
	"github.com/sing3demons/go-mongo-api/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type service struct {
	db     *repository
	logger *zap.Logger
}

func NewService(db *repository, logger *zap.Logger) *service {
	return &service{db: db, logger: logger}
}

func (s *service) Register(account User) (*User, error) {
	exit, err := s.db.FindByEmail(account.Email)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if exit != nil {
		return nil, fmt.Errorf("email already exist")
	}

	account.ID = primitive.NewObjectID()
	hashPassword, err := hashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	account.Password = string(hashPassword)

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

func (s *service) Login(email, password string) (any, error) {
	u, err := s.db.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = comparePassword(password, u.Password)
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
	// hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, err
	// }
	argon := argon2.DefaultConfig()
	hashPassword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}

func comparePassword(pwd, encoded string) error {
	// err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(pwd))
	// if err != nil {
	// 	return  err
	// }
	_, err := argon2.VerifyEncoded([]byte(pwd), []byte(encoded))

	if err != nil {
		logrus.Error("can't compare hash and password err: ", err)
		return err
	}

	return nil
}
