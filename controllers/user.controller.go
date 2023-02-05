package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/siddharth1825/golang-gorm-psql/initializers"
	"github.com/siddharth1825/golang-gorm-psql/models"
)

func CreateUserHandler(c *fiber.Ctx) error {
	var payload *models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"fail","message":err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil{
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newUser := models.User{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(),"duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist, please use another title"})
	} else if result.Error != nil{
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"error","message":result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status":"success","data":fiber.Map{"user":newUser}})
}

func FindUsers(c *fiber.Ctx) error {
	var page = c.Query("page","1")
	var limit = c.Query("limit","10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil{
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"error","message":results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success","results":len(users), "users":users})
}