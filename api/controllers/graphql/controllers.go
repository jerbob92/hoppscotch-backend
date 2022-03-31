package graphql

import (
	"github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
)

func AttachControllers(r *gin.RouterGroup) error {
	r.Any("", graphqlRequest)
	return nil
}

func graphqlRequest(c *gin.Context) {
	reqC := context.GetContext(c)
	reqC.DisableResponses = true

	Handle(reqC, func(req *Request) *graphql.Response {
		c.Set("graphqlC", reqC)
		return Handler.Schema.Exec(c, req.Query, req.OperationName, req.Variables)
	})
}
