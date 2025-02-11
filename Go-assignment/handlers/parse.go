package handlers

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"import-package/config"
	"import-package/models"
)

func ReadAndInsertExcel(filePath string) ([]models.Employee, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel rows: %w", err)
	}

	var employees []models.Employee
	for i, row := range rows {
		if i == 0 || len(row) < 10 {
			continue
		}

		employee := models.Employee{
			FirstName:   row[0],
			LastName:    row[1],
			CompanyName: row[2],
			Address:     row[3],
			City:        row[4],
			County:      row[5],
			Postal:      row[6],
			Phone:       row[7],
			Email:       row[8],
			Web:         row[9],
		}

		employees = append(employees, employee)

		// Insert into MySQL
		err := insertIntoDB(employee)
		if err != nil {
			log.Printf("Error inserting row: %v\n", err)
		}
	}

	// Cache data in Redis
	data, _ := json.Marshal(employees)
	_ = config.SetCache("employees", string(data), 5*time.Minute)

	return employees, nil
}

func GetEmployees(c *gin.Context) {
	cachedData, err := config.GetCache("employees")
	if err == nil {
		c.JSON(200, gin.H{"source": "redis", "data": cachedData})
		return
	}

	rows, err := config.DB.Query("SELECT first_name, last_name, company_name, address, city, county, postal, phone, email, web FROM employees")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(&emp.FirstName, &emp.LastName, &emp.CompanyName, &emp.Address, &emp.City, &emp.County, &emp.Postal, &emp.Phone, &emp.Email, &emp.Web)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}

	data, _ := json.Marshal(employees)
	_ = config.SetCache("employees", string(data), 5*time.Minute)

	c.JSON(200, gin.H{"source": "mysql", "data": employees})
}
