package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramstime/studentrecords/models"
	log "github.com/sirupsen/logrus"
)

// GET /students
// GET /students?id=2
// Find all students
func FindStudents(c *gin.Context) {
	var students []models.Student

	id := c.Query("id")
	if id != "" {
		var student models.Student
		fmt.Printf("id :%v\n", id)
		fmt.Printf("got find request id=%v\n", c.Query("id"))
		if err := models.DB.Where("id = ?", c.Query("id")).First(&student).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			fmt.Printf("error:%v\n", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": student})
		return
	}

	fmt.Printf("got find all request id=%v\n", c.Query("id"))
	models.DB.Find(&students)

	c.JSON(http.StatusOK, gin.H{"data": students})
}

// GET /students/:id
// Find a student
func FindStudent(c *gin.Context) {
	// Get model if exist
	var student models.Student
	id := c.Query("id")
	fmt.Printf("got find request id=%v\n", c.Params.ByName("id"))
	if err := models.DB.Where("id = ?", c.Query("id")).First(&student).Error; err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		var students []models.Student
		fmt.Printf("error:%v\n", err)

		fmt.Printf("got find all request id=%v\n", id)
		models.DB.Find(&students)
		c.JSON(http.StatusOK, gin.H{"data": students})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}

// POST /students
// Create new student
func CreateStudent(c *gin.Context) {
	// Validate input
	var input []models.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var students []models.Student
	// Create student
	for _, data := range input {
		student := models.Student{Title: data.Title, Name: data.Name, Branch: data.Branch,
			Address: data.Address} //Subjects: data.Subjects,
		models.DB.Create(&student)
		if models.DB.Error != nil {
			log.Errorf("got error on insert to DB err:%v", models.DB.Error.Error())
			c.JSON(http.StatusConflict, gin.H{"error": models.DB.Error.Error()})
			return
		}
		students = append(students, student)
	}

	c.JSON(http.StatusOK, gin.H{"data": students})
}

// PATCH /students/:id
// Update a student
func UpdateStudent(c *gin.Context) {
	// Get model if exist
	var student models.Student
	fmt.Printf("got update request id=%v", c.Query("id"))
	if err := models.DB.Where("id = ?", c.Query("id")).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input models.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&student).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": student})
}

// DELETE /students/:id
// Delete a student
func DeleteStudent(c *gin.Context) {
	// Get model if exist
	var student models.Student
	if err := models.DB.Where("id = ?", c.Query("id")).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&student)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
