package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kasasunil/auth-rest-api/internal/utils"
	"log"
	"time"
)

var jwtKey = []byte("secret_auth_rest_api_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// CreateJWTToken generates a new JWT token
func CreateJWTToken(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // Token will expire after 5 minutes (for testing purpose I have put only 5 minutes)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyJWTToken verifies the JWT token from the request
func VerifyJWTToken(ctx *gin.Context) (*Claims, error) {
	tokenStr, err := utils.GetTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Println("ERROR PARSING TOKEN: ", err.Error())
		return nil, err
	}
	if !token.Valid {
		log.Println("SESSION_EXPIRED")
		return nil, fmt.Errorf("SESSION_EXPIRED")
	}
	return claims, nil
}
