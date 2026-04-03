package database

import (
	"hoodhire-chat/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB{
	dsn:=os.Getenv("DB_URL")
	db,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("unable to connect to database:",err)
	}
	r:=db.AutoMigrate(&models.Message{})
	if r!=nil{
		log.Fatal("failed to migrate:",err)
	}
	log.Println("successfully connected to database")
	DB=db
	return db
}
