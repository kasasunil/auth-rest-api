package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kasasunil/auth-rest-api/internal/controllers"
	"github.com/kasasunil/auth-rest-api/internal/entities/revoked_tokens"
	"github.com/kasasunil/auth-rest-api/internal/middlewares"
)

func InitializeRoutes(con *controllers.Controller, rtm *revoked_tokens.RevokedToken) *gin.Engine {
	router := gin.New()
	setupRoutes(router, con, rtm)
	return router
}

func setupRoutes(router *gin.Engine, con *controllers.Controller, rtm *revoked_tokens.RevokedToken) {
	publicGroup := router.Group("public")
	publicGroup.POST("/signup", con.Signup)
	publicGroup.POST("/signin", con.Signin)
	publicGroup.POST("/revoke_token", con.RevokeToken)

	privateGroup := router.Group("private")
	privateGroup.Use(middlewares.VerifyUserSession(rtm))
	privateGroup.GET("/user", con.GetUser)
	privateGroup.GET("/refresh_token", con.RefreshToken)
}
