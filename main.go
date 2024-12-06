package main

import (
	"github.com/kasasunil/auth-rest-api/internal/controllers"
	"github.com/kasasunil/auth-rest-api/internal/database"
	"github.com/kasasunil/auth-rest-api/internal/entities/revoked_tokens"
	"github.com/kasasunil/auth-rest-api/internal/entities/user"
	"github.com/kasasunil/auth-rest-api/internal/router"
	"log"
	"net/http"
)

func main() {
	// Connect to db
	db := database.New()
	err := db.Connect()
	if err != nil {
		log.Println("Failed to connect to the database: ", err)
	} else {
		log.Println("Connected to the database successfully!!!!")
	}

	// Run auto migrations --> I am running migrations to minimize run commands of the application.
	err = db.Db.AutoMigrate(&user.UserModel{}, &revoked_tokens.RevokedTokenModel{})
	if err != nil {
		log.Println("Failed to run auto migrations: ", err)
	} else {
		log.Println("Auto migrations run successfully!!!!")
	}

	userModel := user.New(db)
	revokedTokenModel := revoked_tokens.New(db)

	// Initialize controller
	cont := controllers.New(userModel, revokedTokenModel)

	// Initialize router
	ginHandler := router.InitializeRoutes(cont, revokedTokenModel)

	// Run the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: ginHandler,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Failed to run the server: ", err)
	} else {
		log.Println("Server is running successfully!!!!")
	}
}
