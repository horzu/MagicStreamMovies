package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var secretKey string = os.Getenv("SECRET_KEY")
var refreshSecretKey string = os.Getenv("SECRET_REFRESH_KEY")
var userCollection *mongo.Collection = database.OpenCollection("users")

func GenerateTokens(email, firstName, lastName, role, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(refreshSecretKey))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"token":         signedToken,
		"refresh_token": signedRefreshToken,
		"updated_at":    updateAt,
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": updateData})
	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		return "", errors.New("bearer token missing")
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*SignedDetails); !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	return claims, nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exist := c.Get("userId")
	if !exist {
		return "", errors.New("userId does not exist in this context")
	}

	id, ok := userId.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil
}
