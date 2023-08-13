package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	FirstName string             `bson:"first_name,omitempty" json:"first_name"`
	LastName  string             `bson:"last_name,omitempty" json:"last_name"`
	Password  string             `bson:"password" json:"-"`
	Email     string             `bson:"email" json:"email"`
}

type Login struct {
	Password  string             `bson:"password" binding:"required"`
	Email     string             `bson:"email" binding:"required"`
}
