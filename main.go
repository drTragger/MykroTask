package main

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/drTragger/MykroTask/api/routers"
	"github.com/drTragger/MykroTask/config"
	"github.com/drTragger/MykroTask/repositories"
	"github.com/drTragger/MykroTask/services"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Initialize the database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)

	// Set up router
	router := routers.SetupRouter(userController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
