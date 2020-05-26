package handlers

import (
	"net/http"
	"todo-api/models"

	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID   uint   `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type TaskHandler struct {
	// tg *models.TaskGorm
	tg models.TaskService
}

// func NewTaskHandler(tg *models.TaskGorm) *TaskHandler {
// 	return &TaskHandler{tg}
// }
func NewTaskHandler(tg models.TaskService) *TaskHandler {
	return &TaskHandler{tg}
}

func (th *TaskHandler) ListTask(c *gin.Context) {
	tts, err := th.tg.ListTask()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	tasks := []Task{}
	for _, tt := range tts {
		tasks = append(tasks, Task{
			ID:   tt.ID,
			Task: tt.Task,
			Done: tt.Done,
		})
	}

	c.JSON(http.StatusOK, tasks)
}

type NewTask struct {
	Task string
}

func (th *TaskHandler) CreateTask(c *gin.Context) {
	task := new(NewTask)
	if err := c.BindJSON(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	taskTable := new(models.TaskTable)
	taskTable.Task = task.Task
	taskTable.Done = false
	if err := th.tg.CreateTask(taskTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, Task{
		ID:   taskTable.ID,
		Task: taskTable.Task,
		Done: taskTable.Done,
	})
}

func (th *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}

	if err := th.tg.DeleteTask(uint(id)); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Status(204)
}

type UpdateTask struct {
	Done bool
}

func (th *TaskHandler) UpdateTaskFn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}

	task := new(UpdateTask)
	if err := c.BindJSON(task); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	found, err := th.tg.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	found.Done = task.Done
	if err := th.tg.UpdateTask(found); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Status(204)
}
