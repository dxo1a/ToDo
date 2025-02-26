package handlers

import (
	"ToDo/config"
	"ToDo/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetProfile godoc
//
//	@Summary		Профиль
//	@Description	Получение профиля пользователя
//	@Tags			profile
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} 						models.User
//	@Failure		401	"Invalid token claims"
//	@Failure		404	{object} models.ErrorResponse	"User not found"
//	@Router			/profile [get]
func GetProfile(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateProfile godoc
//
//	@Summary		Изменение профиля
//	@Description	тут описание
//	@Tags			profile
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.ProfileForSwagger	true	"JSON объект с данными пользователя"
//	@Success		200	{object} 						models.Task
//	@Failure		401	"Invalid token claims"
//	@Failure		400	{object} models.ErrorResponse	"Invalid input"
//	@Failure		404	{object} models.ErrorResponse	"User not found"
//	@Failure		500	{object} models.ErrorResponse	"Failed to update profile"
//	@Router			/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	userID := uint(claims["user_id"].(float64))

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"message": err.Error(),
		})
	}

	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Avatar = updatedUser.Avatar

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update profile",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
