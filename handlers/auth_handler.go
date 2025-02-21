package handlers

import (
	"ToDo/config"
	"ToDo/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Register godoc
//
//	@Summary		Регистрация
//	@Description	регистрация по username и password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UserDataRequest	true	"JSON объект с данными пользователя"
//	@Success		201	{object} models.Response "User registered"
//	@Failure		400	{object} models.ErrorResponse "Invalid input"
//	@Failure		409	{object} models.ErrorResponse "User already exists"
//	@Failure		500	{object} models.ErrorResponse "Could not hash password"
//	@Router			/register [post]
func Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	var existingUser models.User
	if err := config.DB.Where("username =?", input.Username).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{
			"error":   "User already exists",
			"message": "Пользователь с таким именем уже существует",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Could not hash password",
			"message": err.Error(),
		})
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Could not create user",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered"})
}

// Login godoc
//
//	@Summary		Авторизация
//	@Description	авторизация по username и password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UserDataRequest	true	"JSON объект с данными пользователя"
//	@Success		200	{object} models.Response "Login successful"
//	@Failure		400	{object} models.ErrorResponse "Invalid input"
//	@Failure		401	{object} models.ErrorResponse "Invalid credentials"
//	@Failure		500	{object} models.ErrorResponse "Could not generate token"
//	@Router			/login [post]
func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"message": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"message": err.Error(),
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Could not generate token",
			"message": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	return c.Status(200).JSON(fiber.Map{"message": "Login successful"})
}

// Logout godoc
//
// @Summary Выход из учётной записи пользователя
// @Description тут типа должно быть описание
// @Tags auth
// @Accept json
// @Produce json
// @Success	200	{object} models.Response "Logged out"
// @Router /logout [post]
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})
	return c.Status(200).JSON(fiber.Map{"message": "Logged out"})
}
