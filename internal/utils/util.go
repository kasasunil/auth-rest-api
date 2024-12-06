package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func GetTokenFromHeader(ctx *gin.Context) (string, error) {
	tokenStr := ctx.GetHeader("Authorization")
	log.Println("Token: ", tokenStr)
	if tokenStr == "" {
		return "", fmt.Errorf("TOKEN_NOT_FOUND")
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	return tokenStr, nil
}
