package api

import (
	"github.com/jerbob92/hoppscotch-backend/api/controllers"
	"github.com/jerbob92/hoppscotch-backend/db"
	"github.com/jerbob92/hoppscotch-backend/helpers/responses"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ginlogrus "github.com/toorop/gin-logrus"
)

func StartAPI() error {
	environment := viper.GetString("api.environment")
	switch environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "development":
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	if environment != "development" {
		log.SetFormatter(&log.JSONFormatter{})
		r.Use(ginlogrus.Logger(log.StandardLogger()), gin.CustomRecovery(responses.RecoveryHandler))
	} else {
		r.Use(gin.Logger())
	}

	// Automatic panic recovery
	r.Use(gin.Recovery())

	// Configuration of CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	if environment == "development" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = viper.GetStringSlice("allowed_domains")
	}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConfig))

	// Attach a DB session to every request.
	r.Use(db.AttachRequestSession())

	// Bind all the routing handlers
	if err := controllers.AttachControllers(r); err != nil {
		return err
	}

	if viper.GetBool("api.ssl.enabled") {
		return r.RunTLS(":"+viper.GetString("api.port"), viper.GetString("api.ssl.certificate"), viper.GetString("api.ssl.key"))
	}

	return r.Run(":" + viper.GetString("api.port"))
}
