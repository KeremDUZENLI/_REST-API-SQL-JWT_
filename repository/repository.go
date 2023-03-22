package repository

import (
	"context"

	"jwt-project/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct{}

type MongoRepository interface {
	AddUser(c *gin.Context, ctx context.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error)
	GetResults(c *gin.Context, ctx context.Context) *mongo.Cursor
}

func NewRepository() MongoRepository {
	return &mongoRepository{}
}
