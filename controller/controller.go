package controller

import (
	"context"
	"jwt-project/database"
	"jwt-project/helper"
	"jwt-project/models"
	"jwt-project/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp() gin.HandlerFunc {
	return service.InsertInDatabase
}

func Login() gin.HandlerFunc {
	return service.FindInDatabase
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

		err := database.Collection(database.MongoClient, "table").FindOne(ctx, bson.M{"userid": personId}).Decode(&person)
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

		result, _ := database.Collection(database.MongoClient, "table").Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		var allUsers []bson.M
		result.All(ctx, &allUsers)

		c.JSON(http.StatusOK, allUsers)
	}
}
