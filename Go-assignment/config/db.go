package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// ConnectDB initializes MySQL connection
func ConnectDB() {
	var err error
	dsn := "root:testdb123456789@tcp(127.0.0.1:3306)/storedb?charset=utf8mb4&parseTime=true&loc=Local"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}

	// Testing connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging MySQL: %v", err)
	}

	fmt.Println("Connected to MySQL (Existing Database)")

	//Ensure table exists
	ensureTableExists()
}

// ensureTableExists checks if the employees table exists and creates
func ensureTableExists() {
	query := `CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		company_name VARCHAR(100),
		address VARCHAR(255),
		city VARCHAR(50),
		county VARCHAR(50),
		postal VARCHAR(20),
		phone VARCHAR(20),
		email VARCHAR(100),
		web VARCHAR(100)
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error ensuring table exists: %v", err)
	}

	fmt.Println("Table ensured in MySQL (employees)")
}
