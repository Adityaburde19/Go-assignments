package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"import-package/config"
	"import-package/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func InsertEmployees(c *gin.Context) {
	var employees []models.Employee

	if err := c.ShouldBindJSON(&employees); err != nil {
		var singleEmployee models.Employee
		if err := c.ShouldBindJSON(&singleEmployee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}
		employees = []models.Employee{singleEmployee} // Store as a single-element array
	}

	var insertedRecords []models.Employee
	var skippedRecords []models.Employee

	for _, emp := range employees {
		// Check if the record already exists
		exists, err := recordExists(emp)
		if err != nil {
			log.Printf("Error checking record existence: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if exists {
			skippedRecords = append(skippedRecords, emp) // Store skipped records
		} else {
			// Insert only if the record does not exist
			err := insertIntoDB(emp)
			if err != nil {
				log.Printf("Error inserting record: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert record"})
				return
			}
			insertedRecords = append(insertedRecords, emp) // Store inserted records
		}
	}

	// Return a response indicating which records were inserted or skipped
	c.JSON(http.StatusOK, gin.H{
		"inserted_records": insertedRecords,
		"skipped_records":  skippedRecords,
	})
}

func recordExists(emp models.Employee) (bool, error) {
	query := `SELECT COUNT(*) FROM employees WHERE first_name = ? AND last_name = ? AND email = ?`
	var count int

	err := config.DB.QueryRow(query, emp.FirstName, emp.LastName, emp.Email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil // If count > 0, record exists
}

func insertIntoDB(emp models.Employee) error {
	query := `INSERT INTO employees (first_name, last_name, company_name, address, city, county, postal, phone, email, web) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := config.DB.Exec(query, emp.FirstName, emp.LastName, emp.CompanyName, emp.Address, emp.City, emp.County, emp.Postal, emp.Phone, emp.Email, emp.Web)
	return err
}

func UpdateEmployee(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	employeeID, err := strconv.Atoi(idStr)
	if err != nil || employeeID <= 0 {
		log.Printf("Invalid employee ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Checks if the employee exists
	existingEmployee, err := getEmployeeByID(employeeID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	} else if err != nil {
		log.Printf("Error fetching employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mergedEmployee := mergeUpdates(existingEmployee, updates)

	err = updateInDB(employeeID, mergedEmployee, updates)
	if err != nil {
		log.Printf("Error updating employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

func getEmployeeByID(id int) (models.Employee, error) {
	var emp models.Employee
	query := `SELECT first_name, last_name, company_name, address, city, county, postal, phone, email, web 
	          FROM employees WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&emp.FirstName, &emp.LastName, &emp.CompanyName, &emp.Address,
		&emp.City, &emp.County, &emp.Postal, &emp.Phone, &emp.Email, &emp.Web,
	)
	return emp, err
}

func mergeUpdates(existing models.Employee, updates map[string]interface{}) models.Employee {
	if val, ok := updates["first_name"].(string); ok {
		existing.FirstName = val
	}
	if val, ok := updates["last_name"].(string); ok {
		existing.LastName = val
	}
	if val, ok := updates["company_name"].(string); ok {
		existing.CompanyName = val
	}
	if val, ok := updates["address"].(string); ok {
		existing.Address = val
	}
	if val, ok := updates["city"].(string); ok {
		existing.City = val
	}
	if val, ok := updates["county"].(string); ok {
		existing.County = val
	}
	if val, ok := updates["postal"].(string); ok {
		existing.Postal = val
	}
	if val, ok := updates["phone"].(string); ok {
		existing.Phone = val
	}
	if val, ok := updates["email"].(string); ok {
		existing.Email = val
	}
	if val, ok := updates["web"].(string); ok {
		existing.Web = val
	}

	return existing
}

// updateInDB updates only the provided fields in MySQL
func updateInDB(id int, emp models.Employee, updates map[string]interface{}) error {
	query := "UPDATE employees SET "
	var params []interface{}

	if _, ok := updates["first_name"]; ok {
		query += "first_name = ?, "
		params = append(params, emp.FirstName)
	}
	if _, ok := updates["last_name"]; ok {
		query += "last_name = ?, "
		params = append(params, emp.LastName)
	}
	if _, ok := updates["company_name"]; ok {
		query += "company_name = ?, "
		params = append(params, emp.CompanyName)
	}
	if _, ok := updates["address"]; ok {
		query += "address = ?, "
		params = append(params, emp.Address)
	}
	if _, ok := updates["city"]; ok {
		query += "city = ?, "
		params = append(params, emp.City)
	}
	if _, ok := updates["county"]; ok {
		query += "county = ?, "
		params = append(params, emp.County)
	}
	if _, ok := updates["postal"]; ok {
		query += "postal = ?, "
		params = append(params, emp.Postal)
	}
	if _, ok := updates["phone"]; ok {
		query += "phone = ?, "
		params = append(params, emp.Phone)
	}
	if _, ok := updates["email"]; ok {
		query += "email = ?, "
		params = append(params, emp.Email)
	}
	if _, ok := updates["web"]; ok {
		query += "web = ?, "
		params = append(params, emp.Web)
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := config.DB.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}
