package main

import (
	"github.com/drTragger/MykroTask/api/controllers"
	"github.com/drTragger/MykroTask/api/routers"
	"github.com/drTragger/MykroTask/config"
	"github.com/drTragger/MykroTask/repository"
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

	jwtKey := []byte(cfg.JWTSecret)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	projectMemberRepo := repository.NewProjectMemberRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, jwtKey)
	projectService := services.NewProjectService(projectRepo, projectMemberRepo, db)
	projectMemberService := services.NewProjectMemberService(projectMemberRepo)
	taskService := services.NewTaskService(taskRepo, projectMemberRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	projectController := controllers.NewProjectController(projectService)
	projectMemberController := controllers.NewProjectMemberController(projectMemberService)
	taskController := controllers.NewTaskController(taskService)

	// Set up router
	router := routers.SetupRouter(userController, projectController, projectMemberController, taskController, jwtKey)

	log.Fatal(http.ListenAndServe(":8080", router))
}
