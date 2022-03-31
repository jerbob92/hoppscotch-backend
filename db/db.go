package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const ConnectionOptions = "charset=utf8mb4&parseTime=True&loc=Local"

func FormatDSN(username string, password string, address string, database string) string {
	return fmt.Sprintf("%s:%s@%s/%s?%s", username, password, address, database, ConnectionOptions)
}

func GetDatabaseDSN(database string) string {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	address := viper.GetString("database.address")

	return FormatDSN(username, password, address, database)
}

func GetDefaultDSN() string {
	database := viper.GetString("database.database")

	return GetDatabaseDSN(database)
}

func ConnectDB() error {
	db, err := gorm.Open(mysql.Open(GetDefaultDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	if viper.GetBool("database.debug") {
		DB = DB.Debug()
	}

	return nil
}

func AttachRequestSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", DB.Session(&gorm.Session{}))

		// Process request
		c.Next()
	}
}
