package service

import (
	"context"
	"jwt-project/database"
	"jwt-project/helper"
	"jwt-project/models"
	"jwt-project/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertInDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var person models.Person
	defer cancel()

	c.BindJSON(&person)

	if repository.Exist(c, ctx, person) || repository.IsValid(c, person) {
		return
	}

	person.ID = primitive.NewObjectID()
	*person.Password = repository.HashPassword(*person.Password)
	token, refreshToken := helper.GenerateAllTokens(*person.Email, *person.FirstName, *person.LastName, *person.UserType, person.UserId)
	person.Token = &token
	person.RefreshToken = &refreshToken
	person.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	person.UserId = person.ID.Hex()

	insert := repository.InsertNumberInDatabase(c, ctx, person)
	c.JSON(http.StatusOK, insert)
}

func FindInDatabase(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var person models.Person
	var foundPerson models.Person
	defer cancel()

	c.BindJSON(&person)

	if repository.IsValidEmail(c, ctx, person, &foundPerson) || repository.IsValidPassword(c, person, foundPerson) {
		return
	}

	repository.Update(ctx, foundPerson)

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

	if err := helper.CheckPersonType(c, models.ADMIN); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var allUsers []bson.M
	repository.Results(c, ctx).All(ctx, &allUsers)

	c.JSON(http.StatusOK, allUsers)
}
