package controller

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var personCollection *mongo.Collection = database.Collection(database.MongoClient, "table")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func VerifyPassword(password string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(password))
	check := true
	msg := ""

	if err != nil {
		check = false
		msg = "password is incorrect"
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var person models.Person

		c.BindJSON(&person)

		validationErr := validate.Struct(person)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "out of rules"})
			return
		}

		count, _ := personCollection.CountDocuments(ctx, bson.M{"email": person.Email})
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		person.ID = primitive.NewObjectID()

		hashedPassword := HashPassword(*person.Password)
		person.Password = &hashedPassword

		token, refreshToken := helper.GenerateAllTokens(*person.Email, *person.FirstName, *person.LastName, *person.UserType, *&person.UserId)
		person.Token = &token
		person.RefreshToken = &refreshToken

		person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		person.UserId = person.ID.Hex()

		resultInsertionNumber, _ := personCollection.InsertOne(ctx, person)
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var person models.Person
		var foundPerson models.Person

		if err := c.BindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := personCollection.FindOne(ctx, bson.M{"email": person.Email}).Decode(&foundPerson)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not correct"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*person.Password, *foundPerson.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundPerson.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshToken := helper.GenerateAllTokens(*foundPerson.Email, *foundPerson.FirstName, *foundPerson.LastName, *foundPerson.UserType, foundPerson.UserId)
		helper.UpdateAllTokens(token, refreshToken, foundPerson.UserId)
		personCollection.FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundPerson)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		personId := c.Param("userid")

		if err := helper.MatchPersonTypeToUid(c, personId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var person models.Person

		err := personCollection.FindOne(ctx, bson.M{"userid": personId}).Decode(&person)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckPersonType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err1 := strconv.Atoi(c.Query("recordPerPage"))
		if err1 != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err2 := strconv.Atoi(c.Query("page"))
		if err2 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, _ = strconv.Atoi(c.Query("startIndex"))

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

		result, _ := personCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		var allUsers []bson.M
		result.All(ctx, &allUsers)

		c.JSON(http.StatusOK, allUsers)
	}
}
