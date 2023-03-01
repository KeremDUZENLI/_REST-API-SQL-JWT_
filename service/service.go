package service

import (
	"context"
	"jwt-project/database"
	"jwt-project/helper"
	"jwt-project/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func exist(c *gin.Context, ctx context.Context, person models.Person) bool {
	if count, _ := database.Collection(database.Database(), models.TABLE).CountDocuments(ctx, bson.M{"email": person.Email}); count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
		return true
	}
	return false
}

func inValid(c *gin.Context, person models.Person) bool {
	if validationErr := validator.New().Struct(person); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "out of rules"})
		return true
	}
	return false
}

func inValidEmail(c *gin.Context, ctx context.Context, person models.Person, foundPerson *models.Person) bool {
	err := database.Collection(database.Database(), models.TABLE).FindOne(ctx, bson.M{"email": person.Email}).Decode(&foundPerson)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not correct"})
		return true
	}
	return false
}

func inValidPassword(c *gin.Context, person models.Person, foundPerson models.Person) bool {
	passwordIsValid, msg := VerifyPassword(*person.Password, *foundPerson.Password)

	if !passwordIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return true
	}
	return false
}

func stages(c *gin.Context) (primitive.D, primitive.D, primitive.D) {
	recordPerPage, err1 := strconv.Atoi(c.Query("recordPerPage"))
	if err1 != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err2 := strconv.Atoi(c.Query("page"))
	if err2 != nil || page < 1 {
		page = 1
	}

	startIndex, _ := strconv.Atoi(c.Query("startIndex"))

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

func HashPassword(password string) string {
	encryptionSize := 14
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), encryptionSize)
	return string(bytes)
}

func VerifyPassword(password string, providedPassword string) (bool, string) {
	if err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(password)); err != nil {
		return false, "password is incorrect"
	}
	return true, ""
}

func InsertInDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var person models.Person
	defer cancel()

	c.BindJSON(&person)

	if exist(c, ctx, person) || inValid(c, person) {
		return
	}

	person.ID = primitive.NewObjectID()

	*person.Password = HashPassword(*person.Password)

	token, refreshToken := helper.GenerateAllTokens(*person.Email, *person.FirstName, *person.LastName, *person.UserType, person.UserId)
	person.Token = &token
	person.RefreshToken = &refreshToken

	person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UserId = person.ID.Hex()

	resultInsertionNumber, _ := database.Collection(database.Database(), models.TABLE).InsertOne(ctx, person)
	c.JSON(http.StatusOK, resultInsertionNumber)
}

func FindInDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var person models.Person
	var foundPerson models.Person
	defer cancel()

	c.BindJSON(&person)

	if inValidEmail(c, ctx, person, &foundPerson) || inValidPassword(c, person, foundPerson) {
		return
	}

	token, refreshToken := helper.GenerateAllTokens(*foundPerson.Email, *foundPerson.FirstName, *foundPerson.LastName, *foundPerson.UserType, foundPerson.UserId)
	helper.UpdateAllTokens(token, refreshToken, foundPerson.UserId)
	database.Collection(database.Database(), models.TABLE).FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson)

	c.JSON(http.StatusOK, &foundPerson)
}

func GetFromDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var person models.Person
	defer cancel()

	personId := c.Param("userid")

	if err := helper.MatchPersonTypeToUid(c, personId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.Collection(database.Database(), models.TABLE).FindOne(ctx, bson.M{"userid": personId}).Decode(&person)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func GetallFromDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := helper.CheckPersonType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	matchStage, groupStage, projectStage := stages(c)
	result, _ := database.Collection(database.Database(), models.TABLE).Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})

	var allUsers []bson.M
	result.All(ctx, &allUsers)

	c.JSON(http.StatusOK, allUsers)
}
