package service

import (
	"context"
	"jwt-project/database"
	"jwt-project/helper"
	"jwt-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var table string = "table"
var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
var person models.Person
var foundPerson models.Person

func exist(c *gin.Context) bool {
	if count, _ := database.Collection(database.MongoClient, table).CountDocuments(ctx, bson.M{"email": person.Email}); count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
		return true
	}
	return false
}

func inValid(c *gin.Context) bool {
	if validationErr := validator.New().Struct(person); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "out of rules"})
		return true
	}
	return false
}

func inValidUser(c *gin.Context) bool {
	if *person.Email == "" || *person.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user"})
		return true
	}
	return false
}

func inValidEmail(c *gin.Context) bool {
	err := database.Collection(database.MongoClient, table).FindOne(ctx, bson.M{"email": person.Email}).Decode(&foundPerson)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not correct"})
		return true
	}
	return false
}

func inValidPassword(c *gin.Context) bool {
	passwordIsValid, msg := VerifyPassword(*person.Password, *foundPerson.Password)

	if !passwordIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return true
	}
	return false
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
	c.BindJSON(&person)

	if exist(c) || inValid(c) {
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

	resultInsertionNumber, _ := database.Collection(database.MongoClient, table).InsertOne(ctx, person)
	c.JSON(http.StatusOK, resultInsertionNumber)
}

func FindInDatabase(c *gin.Context) {
	c.BindJSON(&person)

	if inValidUser(c) || inValidEmail(c) || inValidPassword(c) {
		return
	}

	token, refreshToken := helper.GenerateAllTokens(*foundPerson.Email, *foundPerson.FirstName, *foundPerson.LastName, *foundPerson.UserType, foundPerson.UserId)
	helper.UpdateAllTokens(token, refreshToken, foundPerson.UserId)
	database.Collection(database.MongoClient, table).FindOne(ctx, bson.M{"userid": foundPerson.UserId}).Decode(&foundPerson)

	c.JSON(http.StatusOK, &foundPerson)
}
