package models

import (
	"context"
	"encoding/json"
	"fmt"
	"import-package/config"
	"log"
	"time"
)

type Employee struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"county"`
	Postal      string `json:"postal"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Web         string `json:"web"`
}

func FetchAndCacheEmployees() ([]Employee, error) {
	ctx := context.Background()
	cacheKey := "employees"

	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var employees []Employee
		err := json.Unmarshal([]byte(cachedData), &employees)
		if err != nil {
			log.Println("Error unmarshaling Redis data:", err)
			return nil, err
		}
		fmt.Println("Fetched data from Redis cache")
		return employees, nil
	}

	// Fetch data from MySQL if not found in Redis
	rows, err := config.DB.Query("SELECT first_name, last_name, company_name, address, city, county, postal, phone, email, web FROM employee_data")
	if err != nil {
		log.Println("Error fetching data from MySQL:", err)
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.FirstName, &emp.LastName, &emp.CompanyName, &emp.Address, &emp.City, &emp.County, &emp.Postal, &emp.Phone, &emp.Email, &emp.Web); err != nil {
			log.Println("Error scanning MySQL row:", err)
			return nil, err
		}
		employees = append(employees, emp)
	}

	// Store data in Redis for 5-minute only get flushed
	jsonData, _ := json.Marshal(employees)
	err = config.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		log.Println("Error storing data in Redis:", err)
	}

	fmt.Println("Fetched data from MySQL and stored in Redis")
	return employees, nil
}
