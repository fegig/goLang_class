package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type FieldValue struct {
	Field string
	Value any
}

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	fmt.Println("Connecting to database...", os.Getenv("DBHOST"))
	if os.Getenv("DBUSER") == "" {
		os.Setenv("DBUSER", "root")
	}
	if os.Getenv("DBPASS") == "" {
		os.Setenv("DBPASS", "")
	}
	if os.Getenv("DBHOST") == "" {
		os.Setenv("DBHOST", "localhost")
	}
	if os.Getenv("DBPORT") == "" {
		os.Setenv("DBPORT", "3306")
	}
	if os.Getenv("DBNAME") == "" {
		os.Setenv("DBNAME", "recordings")
	}
	if os.Getenv("DBTYPE") == "" {
		os.Setenv("DBTYPE", "mysql")
	}

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func InsertData(tableName string, fieldValues []FieldValue) (int64, error) {
	query := "INSERT INTO " + tableName + " ("
	values := make([]any, len(fieldValues))

	for i, fieldValue := range fieldValues {
		query += fieldValue.Field
		values[i] = fieldValue.Value
		if i < len(fieldValues)-1 {
			query += ", "
		}
	}

	query += ") VALUES ("
	for i := range fieldValues {
		query += "?"
		if i < len(fieldValues)-1 {
			query += ", "
		}
	}
	query += ")"

	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func SelectData(tableName string, fields []string, conditions []FieldValue, orderBy *string, limit *int, offset *int) (*sql.Rows, error) {
	query := "SELECT "
	for i, field := range fields {
		query += field
		if i < len(fields)-1 {
			query += ", "
		}
	}
	query += " FROM " + tableName

	values := make([]any, 0)
	if len(conditions) > 0 {
		query += " WHERE "
		for i, condition := range conditions {
			query += condition.Field + " = ?"
			values = append(values, condition.Value)
			if i < len(conditions)-1 {
				query += " AND "
			}
		}
	}
	if orderBy != nil {
		query += " ORDER BY " + *orderBy
	}
	if limit != nil {
		query += " LIMIT " + strconv.Itoa(*limit)
	}
	if offset != nil {
		query += " OFFSET " + strconv.Itoa(*offset)
	}
	return db.Query(query, values...)
}

func UpdateData(tableName string, updates []FieldValue, conditions []FieldValue) (int64, error) {
	query := "UPDATE " + tableName + " SET "
	values := make([]any, len(updates)+len(conditions))

	for i, update := range updates {
		query += update.Field + " = ?"
		values[i] = update.Value
		if i < len(updates)-1 {
			query += ", "
		}
	}

	if len(conditions) > 0 {
		query += " WHERE "
		for i, condition := range conditions {
			query += condition.Field + " = ?"
			values[len(updates)+i] = condition.Value
			if i < len(conditions)-1 {
				query += " AND "
			}
		}
	}

	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteData(tableName string, conditions []FieldValue) (int64, error) {
	query := "DELETE FROM " + tableName
	if len(conditions) > 0 {
		query += " WHERE "
		values := make([]any, len(conditions))
		for i, condition := range conditions {
			query += condition.Field + " = ?"
			values[i] = condition.Value
			if i < len(conditions)-1 {
				query += " AND "
			}
		}
		result, err := db.Exec(query, values...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected()
	}
	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
