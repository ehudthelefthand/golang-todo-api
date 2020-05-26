package main

import (
	"log"
	"todo-api/handlers"
	"todo-api/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Cross Origin Resource Sharing

// 1. เพิ่ม axios.post ใน api.js
// 2. call function in React Component ใน App.js
// 3. intercept body ด้วย gin context c.BindJSON()
// 4. logging

func main() {
	db, err := gorm.Open(
		"mysql",
		"root:password@tcp(127.0.0.1:3307)/taskdb?parseTime=true",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true) // dev only!

	if err := db.AutoMigrate(
		&models.TaskTable{},
	).Error; err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	tg := models.NewTaskGorm(db)
	// tg := models.NewTaskMongo(db)

	th := handlers.NewTaskHandler(tg)

	r.GET("/tasks", th.ListTask)

	r.POST("/tasks", th.CreateTask)

	r.DELETE("/tasks/:id", th.DeleteTask)

	r.PATCH("/tasks/:id", th.UpdateTaskFn)

	r.Run()
}
