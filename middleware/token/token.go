package token

import (
	"context"
	"log"
	"time"

	"jwt-project/database"
	"jwt-project/database/model"

	"jwt-project/common/env"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserType  string
	Uid       string
	jwt.StandardClaims
}

func keyFunction(token *jwt.Token) (interface{}, error) {
	return []byte(env.SECRET_KEY), nil
}

func GenerateToken(firstName string, lastName string, email string, userType string, uid string) (token string, refreshToken string, err error) {
	var expiresAt int64 = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	claims := &SignedDetails{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UserType:  userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(env.SECRET_KEY)); err != nil {
		return
	}

	if refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(env.SECRET_KEY)); err != nil {
		return
	}

	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, _ := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		keyFunction,
	)

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "token is invalid"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshtoken", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedat", Value: Updated_at})

	upsert := true
	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.Collection(database.Connect(), model.TABLE).UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}
