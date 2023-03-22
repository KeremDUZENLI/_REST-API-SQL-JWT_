package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DtoLogIn struct {
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"email,required"`

	ID           primitive.ObjectID `bson:"_id"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshtoken"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserId       string             `json:"userid"`
}
