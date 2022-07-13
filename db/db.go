package db

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type DatabaseDSN struct {
	driver            string
	database          string
	username          string
	password          string
	host              string
	port              string
	address           string
	connectionOptions string
}

func (dsn *DatabaseDSN) GetMysqlDSN() string {
	var MysqlConnectionOptions = "charset=utf8mb4&parseTime=True&loc=Local"
	if dsn.connectionOptions != "" {
		MysqlConnectionOptions += "&" + dsn.connectionOptions
	}
	// For backward compatibility
	if dsn.address != "" && dsn.driver == "" {
		return fmt.Sprintf("%s:%s@%s/%s?%s", dsn.username, dsn.password, dsn.address, dsn.database, MysqlConnectionOptions)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dsn.username, dsn.password, dsn.host, dsn.port, dsn.database, MysqlConnectionOptions)
}

func (dsn *DatabaseDSN) GetPostgresDSN() string {
	PostgresConnectionOptions := "TimeZone=" + time.Now().Location().String()
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s %s", dsn.host, dsn.username, dsn.password, dsn.database, dsn.port, PostgresConnectionOptions, dsn.connectionOptions)
}

func (dsn *DatabaseDSN) GetMysql() gorm.Dialector {
	return mysql.Open(dsn.GetMysqlDSN())
}

func (dsn *DatabaseDSN) GetPostgres() gorm.Dialector {
	return postgres.Open(dsn.GetPostgresDSN())
}

func (dsn *DatabaseDSN) GetDialector() (gorm.Dialector, error) {
	// Assuming old version won't have driver field but will have address field
	if dsn.address != "" && dsn.driver == "" {
		return dsn.GetMysql(), nil
	}

	switch dsn.driver {
	case "postgres":
		return dsn.GetPostgres(), nil
	case "mysql":
		return dsn.GetMysql(), nil
	default:
		return nil, errors.New("invalid driver")
	}
}

func ConnectDB() error {
	var connectionData = &DatabaseDSN{
		driver:            viper.GetString("database.driver"),
		database:          viper.GetString("database.database"),
		username:          viper.GetString("database.username"),
		password:          viper.GetString("database.password"),
		host:              viper.GetString("database.host"),
		port:              viper.GetString("database.port"),
		address:           viper.GetString("database.address"),
		connectionOptions: viper.GetString("database.connectionOptions"),
	}

	dialector, err := connectionData.GetDialector()

	if err != nil {
		return err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
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
