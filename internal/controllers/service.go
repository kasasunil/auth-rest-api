package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kasasunil/auth-rest-api/internal/entities/revoked_tokens"
	umodel "github.com/kasasunil/auth-rest-api/internal/entities/user"
	"github.com/kasasunil/auth-rest-api/internal/pkg/jwt"
	"github.com/kasasunil/auth-rest-api/internal/utils"
	"log"
	"net/http"
)

type Controller struct {
	user *umodel.User
	rtm  *revoked_tokens.RevokedToken
}

func New(user *umodel.User, rtm *revoked_tokens.RevokedToken) *Controller {
	return &Controller{
		user: user,
		rtm:  rtm,
	}
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup creates a new user
func (con *Controller) Signup(c *gin.Context) {
	user := &User{}
	err := c.ShouldBindBodyWithJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST"})
		return
	}

	// Transforming request body to usermodel
	model := &umodel.UserModel{
		Email:    user.Email,
		Password: user.Password,
	}

	err = con.user.CreateUser(model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_SIGNING_UP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "USER_CREATED"})
}

// Signin - Authorization of user
func (con *Controller) Signin(c *gin.Context) {
	user := &User{}
	err := c.ShouldBindBodyWithJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST"})
		return
	}

	model, err := con.user.FindUserByEmailId(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "USER_NOT_FOUND"})
		return
	}

	// Authorizing user
	if model.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_CREDENTIALS"})
		return
	}

	// Generate JWT token
	token, err := jwt.CreateJWTToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_GENERATING_TOKEN"})
		return
	}

	// Update JWT token in database
	err = con.user.UpdateUserByEmailId(user.Email, &umodel.UserModel{JwtToken: token})

	c.JSON(http.StatusOK, gin.H{
		"message": "USER_LOGGED_IN_SUCCESSFULLY",
		"token":   token,
		"note1":   "Please use this token in Authorization header(in Bearertoken) for private route: /private/user",
		"note2":   "Token will expire after 5 minutes",
	})
}

// GetUser fetches user details from database with the help of email id fetched from JWT token
func (con *Controller) GetUser(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_FETCHING_USER"})
		return
	}

	user, err := con.user.FindUserByEmailId(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_FETCHING_USER"})
		return
	}

	resp := map[string]interface{}{
		"message":   "USER_FETCHED_SUCCESSFULLY_WITH_JWT_TOKEN",
		"email":     user.Email,
		"jwt_token": user.JwtToken,
	}

	c.JSON(http.StatusOK, gin.H{"user": resp})
}

// RefreshToken refreshes the JWT token
func (con *Controller) RefreshToken(c *gin.Context) {
	emailFromCon, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REFRESHING_TOKEN"})
		return
	}

	var email string
	if val, ok := emailFromCon.(string); ok {
		email = val
	}

	// Create a new token
	token, err := jwt.CreateJWTToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REFRESHING_TOKEN"})
		return
	}

	// Add the current token to revoked tokens list
	tokenStr, err := utils.GetTokenFromHeader(c)
	if err != nil {
		log.Println("ERROR_GETTING_TOKEN_FROM_HEADER: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REFRESHING_TOKEN"})
		return
	}
	err = con.rtm.CreateRevokedToken(&revoked_tokens.RevokedTokenModel{Token: tokenStr})
	if err != nil {
		log.Println("ERROR_ADDING_TOKEN_TO_REVOKED_TOKENS: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REFRESHING_TOKEN"})
		return
	}

	// Update the new token for the corresponding user
	err = con.user.UpdateUserByEmailId(email, &umodel.UserModel{JwtToken: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REFRESHING_TOKEN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "TOKEN_REFRESHED_SUCCESSFULLY",
		"new_token": token,
		"note1":     "this token will expire after 5 minutes",
		"note2":     "If you need to refresh token do it before token expires(i.e. within 5 minutes)",
		"note3":     "If you refresh token the current token you refreshed with will be added to revoked tokens list and will not be valid for any further requests",
	})
}

type RevokeToken struct {
	Token string `json:"token"`
}

func (con *Controller) RevokeToken(c *gin.Context) {
	token := &RevokeToken{}
	err := c.ShouldBindBodyWithJSON(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST"})
		return
	}

	// Transform request body to rt model
	rtPayload := &revoked_tokens.RevokedTokenModel{
		Token: token.Token,
	}

	// Revoke token
	err = con.rtm.CreateRevokedToken(rtPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR_REVOKING_TOKEN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "TOKEN_REVOKED_SUCCESSFULLY",
		"revoked_token": token.Token,
		"note":          "This token will be considered as revoked and will not be valid for any further requests",
	})
}
