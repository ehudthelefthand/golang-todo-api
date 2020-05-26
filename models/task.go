package models

import (
	"github.com/jinzhu/gorm"
)

type TaskTable struct {
	gorm.Model
	Task string
	Done bool
}

func (TaskTable) TableName() string {
	return "tasks"
}

type TaskService interface {
	ListTask() ([]TaskTable, error)
	GetTaskByID(id uint) (*TaskTable, error)
	CreateTask(task *TaskTable) error
	UpdateTask(task *TaskTable) error
	DeleteTask(id uint) error
}

var _ TaskService = &TaskGorm{}

type TaskGorm struct {
	db *gorm.DB
}

func NewTaskGorm(db *gorm.DB) TaskService {
	return &TaskGorm{db}
}

func (tg *TaskGorm) ListTask() ([]TaskTable, error) {
	taskTables := []TaskTable{}
	if err := tg.db.Find(&taskTables).Error; err != nil {
		return nil, err
	}
	return taskTables, nil
}

func (tg *TaskGorm) GetTaskByID(id uint) (*TaskTable, error) {
	tt := new(TaskTable)
	if err := tg.db.First(tt, id).Error; err != nil {
		return nil, err
	}
	return tt, nil
}

func (tg *TaskGorm) CreateTask(task *TaskTable) error {
	return tg.db.Create(task).Error
}

func (tg *TaskGorm) UpdateTask(task *TaskTable) error {
	found := new(TaskTable)
	if err := tg.db.Where("id = ?", task.ID).First(found).Error; err != nil {
		return err
	}
	return tg.db.Model(task).Update("done", task.Done).Error
}

func (tg *TaskGorm) DeleteTask(id uint) error {
	tt := new(TaskTable)
	if err := tg.db.Where("id = ?", id).First(tt).Error; err != nil {
		return err
	}
	return tg.db.Delete(tt).Error
}
