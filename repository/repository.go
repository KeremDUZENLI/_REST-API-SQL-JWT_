package repository

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"jwt-project/database"
	"jwt-project/dto"
	"jwt-project/dto/mapper"
	"jwt-project/middleware"
	"jwt-project/middleware/token"

	"jwt-project/database/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct{}

type MongoRepository interface {
	AddUser(c *gin.Context, ctx context.Context, dSU dto.DtoSignUp) (*mongo.InsertOneResult, error)
	GetStages(c *gin.Context) (primitive.D, primitive.D, primitive.D)
	GetResults(c *gin.Context, ctx context.Context) *mongo.Cursor
}

func NewRepository() MongoRepository {
	return &mongoRepository{}
}

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

func (mongoRepository) GetStages(c *gin.Context) (primitive.D, primitive.D, primitive.D) {
	recordPerPage, errorConvertionRecord := strconv.Atoi(c.Query("recordPerPage"))
	if errorConvertionRecord != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, errorConvertionPage := strconv.Atoi(c.Query("page"))
	if errorConvertionPage != nil || page < 1 {
		page = 1
	}

	startIndex, errorConvertionStartIndex := strconv.Atoi(c.Query("startIndex"))
	if errorConvertionStartIndex != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Provide a valid integer start number"})
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

	return matchStage, groupStage, projectStage
}

func (sR mongoRepository) GetResults(c *gin.Context, ctx context.Context) *mongo.Cursor {
	matchStage, groupStage, projectStage := sR.GetStages(c)
	result, _ := database.Collection(database.Connect(), model.TABLE).Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})

	return result
}
