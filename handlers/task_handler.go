package handlers

import (
	"ToDo/config"
	"ToDo/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

var sanitizer = bluemonday.UGCPolicy()

// Get Tasks godoc
//
//	@Summary		Список задач
//	@Description	Получение списка задач
//	@Tags			CRUD
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}						models.Task
//	@Failure		404	{object} models.ErrorResponse 	"Tasks table not found"
//	@Router			/tasks [get]
func GetTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	if err := config.DB.Select("id", "title", "description", "created_at").Find(&tasks).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Tasks table not found",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(tasks)
}

// Get Task godoc
//
//	@Summary		Задача #
//	@Description	Получение списка задач по ID
//	@Tags			CRUD
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}						models.Task
//	@Failure		404 {object} models.ErrorResponse	"Task not found"
//	@Router			/tasks/{id} [get]
func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := config.DB.First(&task, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Task not found",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(task)
}

// Create Task godoc
//
//		@Summary		Создание задачи
//		@Description	какое-то описание
//		@Tags			CRUD
//		@Accept			json
//		@Produce		json
//		@Param			request	body	models.TaskDataRequest	true	"JSON объект с данными задачи"
//		@Success		201	{object}						models.Task
//		@Failure		400 {object} models.ErrorResponse 	"Invalid input"
//	 	@Failure		500 {object} models.ErrorResponse	"Failed to create task"
//		@Router			/tasks [post]
func CreateTask(c *fiber.Ctx) error {
	var input models.TaskDataRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	task := models.Task{
		Title:       sanitizer.Sanitize(input.Title),
		Description: sanitizer.Sanitize(input.Description),
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create task",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// Delete Task godoc
//
//	@Summary		Удалить задачу
//	@Description	Удаление задачи по ID
//	@Tags			CRUD
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object} models.Response		"Task Deleted"
//	@Failure		404	{object} models.ErrorResponse 	"Task not found"
//	@Router			/tasks/{id} [delete]
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&models.Task{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task deleted"})
}

// Update Task godoc
//
//	@Summary		Изменить задачу
//	@Description	Изменение задачи
//	@Tags			CRUD
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Param			request	body	models.TaskDataRequest	true	"JSON объект с данными задачи"
//	@Success		200	{object} 						models.Task
//	@Failure		404	{object} models.ErrorResponse	"Invalid input
//	@Failure		404	{object} models.ErrorResponse	"Task not found"
//	@Router			/tasks/{id} [put]
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.TaskDataRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	result := config.DB.Model(&models.Task{}).Where("id = ?", id).Updates(models.Task{
		Title:       sanitizer.Sanitize(input.Title),
		Description: sanitizer.Sanitize(input.Description),
	})

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task updated"})
}
