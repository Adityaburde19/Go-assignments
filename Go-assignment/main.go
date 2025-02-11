package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"import-package/config"
	"import-package/handlers"
	"import-package/models"
	"log"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	r := gin.Default()

	// API
	r.GET("/import-excel", func(c *gin.Context) {
		filePath := "C:\\Users\\Aditya\\Downloads\\Sample_Employee_data_xlsx(1)(8).xlsx"
		data, err := handlers.ReadAndInsertExcel(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	})

	// API to fetch employee data
	r.GET("/employees", handlers.GetEmployees)
	r.POST("/insert-employees", handlers.InsertEmployees)
	r.PUT("/update/:id", handlers.UpdateEmployee) // Update existing employee by ID

	employees, err := models.FetchAndCacheEmployees()
	if err != nil {
		log.Fatal("Error fetching employees:", err)
	}

	fmt.Println("Employee Data:", employees)

	r.Run(":8081")
}
