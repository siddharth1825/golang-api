package main

import (
	"context"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	gofiberfirebaseauth "github.com/sacsand/gofiber-firebaseauth"
	"github.com/siddharth1825/golang-gorm-psql/controllers"
	"github.com/siddharth1825/golang-gorm-psql/initializers"
	"google.golang.org/api/option"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load env variables", err)
	}

	initializers.ConnectDB(&config)

}

func main() {
	//loading config
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load env variables", err)
	}

	//firebase
	serviceAccountKeyFilePath, err :=
	filepath.Abs("./serviceAccount.json")
	if err != nil {
		panic("unable to load service account")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	fireapp , err := firebase.NewApp(context.Background(),nil,opt)
	if err!= nil {
		panic("firebase load error")
	}

	//initialize app
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api",micro)

	//middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))
	app.Use(gofiberfirebaseauth.New(gofiberfirebaseauth.Config{
		FirebaseApp: fireapp,
		IgnoreUrls: []string{"GET::/api/users","POST::/api/users","GET::/api/songs","POST::/api/songs"},
	}))

	//routes
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

	//running server
	log.Fatal(app.Listen(":" + config.ServerPort))
		
}
