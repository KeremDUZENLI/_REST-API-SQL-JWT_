package repository

import (
	"context"
	"net/http"
	"strconv"

	"jwt-project/common/constants"
	"jwt-project/database"
	"jwt-project/database/model"
	"jwt-project/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

/*
func VerifyPassword(password string, providedPassword string) (bool, string) {
	if err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(password)); err != nil {
		return false, "password is incorrect"
	}
	return true, constants.EMPTY_STRING
}
*/

func HashPassword(password string) string {
	encryptionSize := 14
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), encryptionSize)
	return string(bytes)
}

func Exist(c *gin.Context, ctx context.Context, person model.Person) bool {
	if count, _ := database.Collection(database.Database(), constants.TABLE).CountDocuments(ctx, bson.M{"email": person.Email}); count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
		return true
	}
	return false
}

func IsValid(c *gin.Context, person model.Person) bool {
	if validationErr := validator.New().Struct(person); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "out of rules"})
		return true
	}
	return false
}

func Update(ctx context.Context, foundPerson model.Person) {
	token, refreshToken := helper.GenerateAllTokens(*foundPerson.Email, *foundPerson.FirstName, *foundPerson.LastName, *foundPerson.UserType, foundPerson.UserId)
	helper.UpdateAllTokens(token, refreshToken, foundPerson.UserId)
	database.Collection(database.Database(), constants.TABLE).FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson)
}

func InsertNumberInDatabase(c *gin.Context, ctx context.Context, person model.Person) *mongo.InsertOneResult {
	resultInsertionNumber, _ := database.Collection(database.Database(), constants.TABLE).InsertOne(ctx, person)
	return resultInsertionNumber
}

func Stages(c *gin.Context) (primitive.D, primitive.D, primitive.D) {
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

	matchStage := bson.D{{"$match", bson.D{{}}}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"_id", "null"}}},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}

	return matchStage, groupStage, projectStage
}

func Results(c *gin.Context, ctx context.Context) *mongo.Cursor {
	matchStage, groupStage, projectStage := Stages(c)
	result, _ := database.Collection(database.Database(), constants.TABLE).Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})

	return result
}
