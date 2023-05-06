package main

import (
	"github.com/ramstime/studentrecords/controllers"
	"github.com/ramstime/studentrecords/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect to database
	models.ConnectDatabase()

	// Routes
	r.GET("/students", controllers.FindStudents)
	r.GET("/students/:id", controllers.FindStudent)
	r.POST("/students", controllers.CreateStudent)
	r.PATCH("/students/", controllers.UpdateStudent)
	r.DELETE("/students/", controllers.DeleteStudent)

	// Run the server
	r.Run()
}
