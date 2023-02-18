package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/siddharth1825/golang-gorm-psql/initializers"
	"github.com/siddharth1825/golang-gorm-psql/models"
	"gorm.io/gorm"
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
		Email: payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(),"duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exist, please use another Email"})
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

func UpdateUsers(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var payload *models.UpdateUserSchema

	if err := c.BodyParser(&payload); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"fail","message":err.Error()})
	}

	var user models.User
	results := initializers.DB.First(&user,"id= ?",userId)
	if err := results.Error; err != nil{
		if err == gorm.ErrDryRunModeUnsupported{
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"fail","message":"No suer with that id"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"fail","message":err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Email != ""{
		updates["Email"]=payload.Email
	}

	updates["updated_at"]=time.Now()

	initializers.DB.Model(&user).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success","data":fiber.Map{"user":user}})
}

func FindUserById(c *fiber.Ctx) error{
	userId := c.Params("userId")

	var user models.User
	results := initializers.DB.Preload("Songs").First(&user,"email= ?",userId)
	if err := results.Error; err != nil{
		if err == gorm.ErrRecordNotFound{
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"fail","message":"No user with that email exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"fail","message":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success","data":fiber.Map{"user":user}})
}

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	result := initializers.DB.Delete(&models.User{},"email= ?",userId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"fail","message":"No user with that email exists"})
	} else if result.Error != nil{
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status":"error","message":result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

