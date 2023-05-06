package models

import (
	"os"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbDSN string //= "mysql://rams:rams@tcp(127.0.0.1:3306)/students"
func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Info("We are getting the env values")
	}
}
func ConnectDatabase() {

	//database, err := gorm.Open("sqlite3", "test.db")
	database := mysqlConnect()
	// if err != nil {
	// 	panic("Failed to connect to database!")
	// }

	database.AutoMigrate(&Student{})

	DB = database
}

func mysqlConnect() (db *gorm.DB) {
	// splitting DSN to only use the string after mysql://
	dbDSN = "mysql://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?parseTime=true"
	splitDSN := strings.Split(dbDSN, "mysql://")
	db, err := gorm.Open(mysql.Open(splitDSN[1]), &gorm.Config{})
	if err != nil {
		log.Errorf("Error Connecting MySQL %s, %s", dbDSN, err.Error())
		panic("Failed to connect to database! " + dbDSN)
	}
	return db
}
