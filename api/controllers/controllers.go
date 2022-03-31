package controllers

import (
	"github.com/jerbob92/hoppscotch-backend/api/controllers/graphql"

	"github.com/gin-gonic/gin"
)

func AttachControllers(engine *gin.Engine) error {
	if err := graphql.AttachControllers(engine.RouterGroup.Group("/graphql")); err != nil {
		return err
	}

	return nil
}
