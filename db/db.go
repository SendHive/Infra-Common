package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IDataBaseService interface {
	InitDB() (*gorm.DB, error)
}

type DataBaseService struct {
	Db *gorm.DB
}

func NewDbRequest() (IDataBaseService, error) {
	return &DataBaseService{}, nil
}

func (db *DataBaseService) InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	var err error
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("the error while creating the database connection: ", err.Error())
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		fmt.Println("Error getting underlying sql.DB: ", err.Error())
		return nil, err
	}

	sqlDB.SetMaxOpenConns(20)                  
	sqlDB.SetMaxIdleConns(5)                   
	sqlDB.SetConnMaxLifetime(30 * time.Minute) 
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  

	db.Db = conn
	return conn, nil
}
