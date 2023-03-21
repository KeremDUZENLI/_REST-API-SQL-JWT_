package repository

import (
	"context"
	"errors"
	"time"

	"jwt-project/database"
	"jwt-project/dto"
	"jwt-project/dto/mapper"
	"jwt-project/middleware"
	"jwt-project/middleware/token"

	"jwt-project/database/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mongoRepository) AddUser(c *gin.Context, ctx context.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error) {
	aMap := mapper.MapperSignUp(&dSU)

	setValues(&aMap)

	resultInsertionNumber, err := database.Collection(database.Connect(), model.TABLE).InsertOne(ctx, aMap)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}
	return resultInsertionNumber, nil
}

func setValues(person *model.Person) error {
	person.ID = primitive.NewObjectID()

	person.Password, _ = middleware.HashPassword(person.Password)
	_, errPassword := middleware.HashPassword(person.Password)

	token, refreshToken, errToken := token.GenerateToken(person.Email, person.FirstName, person.LastName, person.UserType, person.UserId)
	person.Token = token
	person.RefreshToken = refreshToken

	person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	person.UserId = person.ID.Hex()

	if errPassword != nil && errToken != nil {
		return errors.New("error setValues")
	}

	return nil
}
