package db

import (
	"fmt"

	// _ "github.com/lib/pq"
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

var db *gorm.DB
var dbs *sql.DB

var postgresStatmente string
var errors error

const (
	username = "samuael"
	password = "samuaelfirst"
	host     = "localhost"
	dbname   = entity.DBName
)

// InitializPostgres  r
func InitializPostgres() (*gorm.DB, error) {
	// Preparing the statmente
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)
	db, errors = gorm.Open("postgres", postgresStatmente)
	if errors != nil {
		panic(errors)
	}
	fmt.Println("DB Connected Succesfully ")
	return db, nil
}
