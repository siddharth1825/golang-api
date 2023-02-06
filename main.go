package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/siddharth1825/golang-gorm-psql/controllers"
	"github.com/siddharth1825/golang-gorm-psql/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load env variables", err)
	}

	initializers.ConnectDB(&config)

}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load env variables", err)
	}

	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api",micro)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Route("/users",func(router fiber.Router){
		router.Post("/", controllers.CreateUserHandler)
		router.Get("",controllers.FindUsers)
	})

	micro.Route("/songs",func(router fiber.Router){
		router.Post("/", controllers.CreateSongs)
		router.Get("",controllers.FindSongs)
	})

	micro.Route("/users/:userId", func(router fiber.Router) {
		router.Delete("",controllers.DeleteUser)
		router.Get("",controllers.FindUserById)
		router.Patch("",controllers.UpdateUsers)
	})

	log.Fatal(app.Listen(":" + config.ServerPort))
		
}
