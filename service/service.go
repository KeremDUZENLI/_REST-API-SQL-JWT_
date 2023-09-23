package service

import (
	"jwt-project/database/model"
	"jwt-project/dto"
	"jwt-project/repository"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoService struct {
	mongoRepository repository.MongoRepository
}

type MongoService interface {
	CreateUser(c *gin.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error)
	FindUser(c *gin.Context, dLI dto.DtoLogIn) (*model.Person, error)
	GetUserByID(c *gin.Context, dGU dto.GetUser, personId string) (model.Person, error)
	GetAllUsers(c *gin.Context, allUsers []primitive.M) ([]primitive.M, error)
}

func NewService(repository repository.MongoRepository) MongoService {
	return &mongoService{mongoRepository: repository}
}
