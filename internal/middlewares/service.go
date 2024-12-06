package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/kasasunil/auth-rest-api/internal/entities/revoked_tokens"
	"github.com/kasasunil/auth-rest-api/internal/pkg/jwt"
	"github.com/kasasunil/auth-rest-api/internal/utils"
	"net/http"
)

// VerifyUserSession verifies the user session and sets the email in the context
func VerifyUserSession(rtm *revoked_tokens.RevokedToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if the token is present in the revoked tokens list
		isPresent := IsTokenPresentInRevokedTokens(c, rtm)
		if isPresent {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "TOKEN_REVOKED_ALREADY",
				"msg":   "Token is revoked. Please signin again. Revoked tokens are not allowed.",
				"note":  "With revoked tokens, even if they have not expired, user cannot refresh the token.\nFOr new token he should signin again.",
			})
			c.Abort()
			return
		}

		claims, err := jwt.VerifyJWTToken(c)
		if err != nil {
			if c.FullPath() == "/private/refresh_token" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "ERROR_REFRESHING_TOKEN/TOKEN_ALREADY_EXPIRED"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

// IsTokenPresentInRevokedTokens checks if the token is present in the revoked tokens list
func IsTokenPresentInRevokedTokens(c *gin.Context, rtm *revoked_tokens.RevokedToken) bool {
	tokenStr, err := utils.GetTokenFromHeader(c)
	if err != nil {
		return false
	}

	// Find the token in the revoked tokens list
	token, err := rtm.FindRevokedTokenByTokenId(tokenStr)
	if err != nil {
		return false
	}
	if token != nil && token.Token == tokenStr {
		return true
	}

	return false
}
