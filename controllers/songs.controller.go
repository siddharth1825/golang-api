package controllers

import(

	"strconv"
	// "strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/siddharth1825/golang-gorm-psql/initializers"
	"github.com/siddharth1825/golang-gorm-psql/models"

)

func CreateSongs(c *fiber.Ctx) error {
	var payload *models.CreateSongSchema

	if err := c.BodyParser(&payload); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"fail","message":err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil{
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newSong := models.Songs{
		Link: payload.Link,
		UserEmail: payload.UserEmail,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newSong)

	if result.Error != nil{
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"error","message":result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status":"success","data":fiber.Map{"user":newSong}})
}

func FindSongs(c *fiber.Ctx) error{
	var page = c.Query("page","1")
	var limit = c.Query("limit","10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var songs []models.Songs
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&songs)
	if results.Error != nil{
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"error","message":results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success","results":len(songs), "users":songs})
}