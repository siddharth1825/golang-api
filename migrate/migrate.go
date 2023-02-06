package main

import (
	"fmt"
	"log"

	"github.com/siddharth1825/golang-gorm-psql/initializers"
	"github.com/siddharth1825/golang-gorm-psql/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Songs{})
	fmt.Println("? Migration complete")
}

