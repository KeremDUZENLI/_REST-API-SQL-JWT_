package service

import (
	"context"
	"errors"
	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/dto"
	"jwt-project/dto/mapper"
	"jwt-project/middleware/token"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

func (sS mongoService) FindUser(c *gin.Context, dLI dto.DtoLogIn) (*model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foundPerson, _ := dto.Find(ctx, dLI)

	if !foundPerson.IsValidEmail(dLI.Email) || !foundPerson.IsValidPassword(dLI.Password) {
		return &model.Person{}, errors.New("invalid email or password")
	}

	aMap := mapper.MapperLogin(foundPerson)

	sS.update(ctx, aMap)
	return &aMap, nil
}

func (mongoService) update(ctx context.Context, foundPerson model.Person) error {
	firstToken, refreshToken, err := token.GenerateToken(foundPerson.Email, foundPerson.FirstName, foundPerson.LastName, foundPerson.UserType, foundPerson.UserId)
	token.UpdateAllTokens(firstToken, refreshToken, foundPerson.UserId)

	database.Collection(database.Connect(), model.TABLE).FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson)

	if err != nil {
		return err
	}

	return nil
}
