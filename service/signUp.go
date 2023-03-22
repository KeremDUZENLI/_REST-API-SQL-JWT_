package service

import (
	"context"
	"errors"
	"jwt-project/dto"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
)

func (sS mongoService) CreateUser(c *gin.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if dSU.IsExist(c) || !dSU.IsObeyRules() {
		return &mongo.InsertOneResult{}, errors.New("email or password either exist or out of rules")
	}

	res, err := sS.mongoRepository.AddUser(c, ctx, dSU)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}

	return res, nil
}
