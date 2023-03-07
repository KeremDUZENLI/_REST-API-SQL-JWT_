package service

import (
	"context"
	"errors"
	"jwt-project/common/constants"
	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/middleware/auth"
	"jwt-project/middleware/token"
	"jwt-project/repository"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertInDatabase(c *gin.Context, person model.Person) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if !person.IsNotExist(c) || !person.IsObeyRules() {
		return &mongo.InsertOneResult{}, errors.New("invalid email or password")
	}

	person.ID = primitive.NewObjectID()
	*person.Password = repository.HashPassword(*person.Password)
	token, refreshToken := token.GenerateToken(*person.Email, *person.FirstName, *person.LastName, *person.UserType, person.UserId)
	person.Token = &token
	person.RefreshToken = &refreshToken
	person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UserId = person.ID.Hex()

	return repository.InsertNumberInDatabase(c, ctx, person), nil
}

func FindInDatabase(c *gin.Context, person model.Person) (*model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foundPerson := model.Find(ctx, person)
	if !foundPerson.IsValidEmail(*person.Email) || !foundPerson.IsValidPassword(*person.Password) {
		return &model.Person{}, errors.New("invalid email or password")
	}

	repository.Update(ctx, *foundPerson)
	return foundPerson, nil
}

func GetFromDatabase(c *gin.Context, person model.Person, personId string) (model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.MatchPersonTypeToUid(c, personId); err != nil {
		return model.Person{}, err
	}

	if err := database.Collection(database.Connect(), constants.TABLE).FindOne(ctx, bson.M{"userid": personId}).Decode(&person); err != nil {
		return model.Person{}, err
	}

	return person, nil
}

func GetallFromDatabase(c *gin.Context, allUsers []primitive.M) ([]primitive.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.CheckPersonType(c, constants.ADMIN); err != nil {
		return []primitive.M{}, err
	}

	repository.Results(c, ctx).All(ctx, &allUsers)
	return allUsers, nil
}
