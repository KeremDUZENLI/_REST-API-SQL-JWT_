package service

import (
	"context"
	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/dto"
	"jwt-project/dto/mapper"
	"jwt-project/middleware/auth"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mongoService) GetUserByID(c *gin.Context, dGU dto.GetUser, personId string) (model.Person, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.MatchPersonTypeToUid(c, personId); err != nil {
		return model.Person{}, err
	}

	if err := database.Collection(database.Connect(), model.TABLE).FindOne(ctx, bson.M{"userid": personId}).Decode(&dGU); err != nil {
		return model.Person{}, err
	}

	aMap := mapper.MapperGetUser(&dGU)

	return aMap, nil
}

func (sS mongoService) GetAllUsers(c *gin.Context, allUsers []primitive.M) ([]primitive.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.CheckPersonType(c, model.ADMIN); err != nil {
		return []primitive.M{}, err
	}

	sS.mongoRepository.GetResults(c, ctx).All(ctx, &allUsers)
	return allUsers, nil
}
