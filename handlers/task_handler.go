package handlers

import (
	"ToDo/config"
	"ToDo/models"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"time"
)

var p = bluemonday.UGCPolicy()

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
	if err := config.DB.Find(&tasks).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Tasks table not found",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(tasks)
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
	if err := config.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Task not found",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(task)
}

// Create Task godoc
//
//	@Summary		Создание задачи
//	@Description	какое-то описание
//	@Tags			CRUD
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.TaskDataRequest	true	"JSON объект с данными задачи"
//	@Success		201	{object}						models.Task
//	@Failure		400 {object} models.ErrorResponse 	"Invalid input"
//	@Router			/tasks [post]
func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	task.Title = p.Sanitize(task.Title)
	task.Description = p.Sanitize(task.Description)
	task.CreatedAt = time.Now()

	config.DB.Create(&task)

	return c.Status(201).JSON(task)
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
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	config.DB.Delete(&task)
	return c.Status(200).JSON(fiber.Map{"message": "Task deleted"}, id)
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
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Task not found",
			"message": err.Error(),
		})
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid input",
			"message": err.Error(),
		})
	}

	config.DB.Save(&task)
	return c.Status(200).JSON(task)
}
