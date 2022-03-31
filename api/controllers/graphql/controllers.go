package graphql

import (
	goctx "context"
	"net/http"

	"github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

func AttachControllers(r *gin.RouterGroup) error {
	r.Any("", graphqlRequest())
	r.Any("ws", graphqlRequest())
	return nil
}

type contextGenerator struct {
}

func (t contextGenerator) BuildContext(ctx goctx.Context, r *http.Request) (goctx.Context, error) {
	c := r.Context().Value("ginctx").(*gin.Context)
	reqC := context.GetContext(c)
	reqC.DisableResponses = true
	return goctx.WithValue(ctx, "graphqlC", reqC), nil
}

func graphqlRequest() gin.HandlerFunc {
	graphQLHandler := graphqlws.NewHandlerFunc(Handler.Schema, Handler, graphqlws.WithContextGenerator(&contextGenerator{}))

	return func(c *gin.Context) {
		clonedRequest := c.Request.WithContext(goctx.WithValue(c.Request.Context(), "ginctx", c))
		graphQLHandler.ServeHTTP(c.Writer, clonedRequest)
	}
}
