package tokens

import (
	"context"
	"ecommerce/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")
var UserData *mongo.Collection = database.UserData(database.Client, "Users")

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}


func (s *SignedDetails) TokenGenerator() (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     s.Email,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Uid:       s.Uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SigningString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SigningString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshtoken, err

}


func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {


	token, err := jwt.ParseWithClaims(signedtoken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Unix() {
		msg = "token is already expired"
		return
	}
	return claims, msg
}
func UpdateToken(signedtoken string , signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	var updateobj primitive.D

	updateobj =append(updateobj, bson.E{Key: "token", Value:signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value:signedrefreshtoken})
	updated_at, _ :=time.Parse((time.RFC3339, time.Now().Format(time.RFC3339)))
	updateobj = append(updateobj,bson.E{Key: "updatedat",Value: updated_at})

	upsert :=true
	filter:=bson.M{"user_id":userid}
	opt := options.UpdateOptions{
		Upsert:&upsert,
	}
	_, err:=UserData.UpdateOne(ctx, filter,bson.D{
		{Key: "$set",Value: updateobj},
	},
	&opt)
	defer cancel()
	if err!=nil{
		log.Panic(err)
		return
	}
}
