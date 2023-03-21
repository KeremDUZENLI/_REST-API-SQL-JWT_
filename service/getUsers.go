package service

import (
	"context"
	"jwt-project/database/model"
	"jwt-project/middleware/auth"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (sS mongoService) GetAllUsers(c *gin.Context, allUsers []primitive.M) ([]primitive.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := auth.CheckPersonType(c, model.ADMIN); err != nil {
		return []primitive.M{}, err
	}

	sS.mongoRepository.GetResults(c, ctx).All(ctx, &allUsers)
	return allUsers, nil
}
