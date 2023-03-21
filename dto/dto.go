package dto

import (
	"context"

	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/middleware"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DtoSignUp) IsExist(ctx context.Context) bool {
	aMap := mapperSignUpLogin(d)
	_, err := Find(ctx, aMap)
	return err == nil
}

func (d DtoSignUp) IsObeyRules() bool { return Validator(d) == nil }

func (d DtoLogIn) IsValidEmail(email string) bool { return email == d.Email }

func (d DtoLogIn) IsValidPassword(password string) bool {
	return middleware.VerifyPassword(password, d.Password)
}

func Find(ctx context.Context, d DtoLogIn) (*DtoLogIn, error) {
	var foundPerson DtoLogIn
	filter := bson.M{"email": d.Email}

	if err := database.Collection(database.Connect(), model.TABLE).
		FindOne(ctx, filter).Decode(&foundPerson); err != nil {
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
