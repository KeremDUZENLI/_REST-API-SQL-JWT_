package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jwt-project/common/constants"
	"jwt-project/database"

	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstname" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastname" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	UserType     *string            `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"createdat"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserId       string             `json:"userid"`
}

func (p Person) IsValidEmail(email string) bool { return email == *p.Email }

func (p Person) IsValidPassword(password string) bool { return Verify(password, *p.Password) }

func (p Person) IsNotExist(ctx context.Context) bool { return !Amount(ctx, p) }

func (p Person) IsObeyRules() bool { return Validator(p) == nil }

func Find(ctx context.Context, person Person) *Person {
	var foundPerson Person
	if err := database.Collection(database.Connect(), constants.TABLE).FindOne(ctx, bson.M{"email": person.Email}).Decode(&foundPerson); err != nil {
		return &person
	}
	return &foundPerson
}

func Verify(password string, providedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func Amount(ctx context.Context, person Person) bool {
	number, err := database.Collection(database.Connect(), constants.TABLE).CountDocuments(ctx, bson.M{"email": person.Email})
	if err != nil {
		return false
	}

	return number > 0
}

func Validator(person Person) error {
	return validator.New().Struct(person)
}
