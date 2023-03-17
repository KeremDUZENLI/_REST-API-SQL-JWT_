package dto

import (
	"context"
	"time"

	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/middleware"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUp----------------------------------------------------------------

type DtoSignUp struct {
	ID           primitive.ObjectID `bson:"_id"`
	Password     string             `json:"password" validate:"required,min=6"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"createdat"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserId       string             `json:"userid"`

	FirstName string `json:"firstname" validate:"required,min=2,max=100"`
	LastName  string `json:"lastname" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"email,required"`
	UserType  string `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
}

func (d DtoSignUp) IsNotExist(ctx context.Context) bool {
	aMap := mapperSignUpLogin(d)
	_, err := Find(ctx, aMap)
	return err != nil
}

func (d DtoSignUp) IsObeyRules() bool { return Validator(d) == nil }

func Find(ctx context.Context, d DtoLogIn) (*DtoLogIn, error) {
	var foundPerson DtoLogIn
	if err := database.Collection(database.Connect(), model.TABLE).FindOne(ctx, bson.M{"email": d.Email}).Decode(&foundPerson); err != nil {
		return &d, err
	}
	return &foundPerson, nil
}

func Validator(d DtoSignUp) error {
	return validator.New().Struct(d)
}

func mapperSignUpLogin(d DtoSignUp) DtoLogIn {
	return DtoLogIn{
		ID:           d.ID,
		Password:     d.Password,
		Token:        d.Token,
		RefreshToken: d.RefreshToken,
		UpdatedAt:    d.UpdatedAt,
		UserId:       d.UserId,

		Email: d.Email,
	}
}

// LogIn----------------------------------------------------------------

type DtoLogIn struct {
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"email,required"`

	ID           primitive.ObjectID `bson:"_id"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshtoken"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserId       string             `json:"userid"`
}

func (d DtoLogIn) IsValidEmail(email string) bool { return email == d.Email }

func (d DtoLogIn) IsValidPassword(password string) bool {
	return middleware.VerifyPassword(password, d.Password)
}

// GetUser----------------------------------------------------------------

type GetUser struct {
	UserId string `json:"userid"`

	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"firstname" validate:"required,min=2,max=100"`
	LastName     string             `json:"lastname" validate:"required,min=2,max=100"`
	Password     string             `json:"password" validate:"required,min=6"`
	Email        string             `json:"email" validate:"email,required"`
	UserType     string             `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"createdat"`
	UpdatedAt    time.Time          `json:"updatedat"`
}
