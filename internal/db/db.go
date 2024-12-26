package db

import (
	"github.com/nelsonfrank/finance-tracker/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// func init() {
// 	// Initialize database connection
// 	dsn := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5438/finance-tracker?sslmode=disable")
// 	var err error
// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

//		// Auto migrate the schema
//		db.AutoMigrate(&model.User{})
//	}
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	db.AutoMigrate(&store.User{})

	return db, nil
}
