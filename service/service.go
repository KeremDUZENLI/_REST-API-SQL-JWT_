package service

import (
	"context"
	"errors"
	"fmt"
	"jwt-project/common/constants"
	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/dto"
	"jwt-project/dto/mapper"
	"jwt-project/middleware/auth"
	"jwt-project/middleware/token"
	"jwt-project/repository"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setValues(person *model.Person) {
	person.ID = primitive.NewObjectID()
	person.Password = repository.HashPassword(person.Password)
	token, refreshToken := token.GenerateToken(person.Email, person.FirstName, person.LastName, person.UserType, person.UserId)
	person.Token = token
	person.RefreshToken = refreshToken
	person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UserId = person.ID.Hex()
}

func InsertInDatabase(c *gin.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if !dSU.IsNotExist(c) || !dSU.IsObeyRules() {
		return &mongo.InsertOneResult{}, errors.New("invalid email or password")
	}

	aMap := mapper.MapperSignUp(&dSU)
	setValues(&aMap)

	return repository.InsertNumberInDatabase(c, ctx, &aMap), nil
}

// ----------------------------------------------------------------

func update(ctx context.Context, foundPerson model.Person) error {
	firstToken, refreshToken := token.GenerateToken(foundPerson.Email, foundPerson.FirstName, foundPerson.LastName, foundPerson.UserType, foundPerson.UserId)
	token.UpdateAllTokens(firstToken, refreshToken, foundPerson.UserId)

	if err := database.Collection(database.Connect(), constants.TABLE).FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func FindInDatabase(c *gin.Context, dLI dto.DtoLogIn) (*model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foundPerson := dto.Find(ctx, dLI)
	if !foundPerson.IsValidEmail(dLI.Email) || !foundPerson.IsValidPassword(dLI.Password) {
		return &model.Person{}, errors.New("invalid email or password")
	}

	aMap := mapper.MapperLogin(foundPerson)

	update(ctx, aMap)
	return &aMap, nil
}

// ----------------------------------------------------------------

func GetFromDatabase(c *gin.Context, dGU dto.GetUser, personId string) (model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.MatchPersonTypeToUid(c, personId); err != nil {
		return model.Person{}, err
	}

	if err := database.Collection(database.Connect(), constants.TABLE).FindOne(ctx, bson.M{"userid": personId}).Decode(&dGU); err != nil {
		return model.Person{}, err
	}

	aMap := mapper.MapperGetUser(&dGU)

	fmt.Println(dGU.UserId)
	fmt.Println(aMap)

	return aMap, nil
}

// ----------------------------------------------------------------

func GetallFromDatabase(c *gin.Context, allUsers []primitive.M) ([]primitive.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.CheckPersonType(c, constants.ADMIN); err != nil {
		return []primitive.M{}, err
	}

	repository.Results(c, ctx).All(ctx, &allUsers)
	return allUsers, nil
}
