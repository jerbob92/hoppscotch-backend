package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/resolvers"
	"github.com/jerbob92/hoppscotch-backend/schema"

	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/relay"
)

var Handler *relay.Handler

func init() {
	s, err := schema.String()
	if err != nil {
		log.Fatalf("reading embedded schema contents: %s", err)
	}

	parsedSchema := graphql.MustParseSchema(
		s,
		&resolvers.BaseQuery{},
		graphql.UseFieldResolvers(),
		graphql.MaxParallelism(5),
		graphql.UseStringDescriptions(),
		graphql.MaxDepth(20), // Just to be sure
		graphql.Logger(graphqlLogger{}),
		graphql.PanicHandler(PanicHandler{}),
	)
	Handler = &relay.Handler{Schema: parsedSchema}
}

type graphqlLogger struct{}

func (g graphqlLogger) LogPanic(ctx context.Context, value interface{}) {
	errorstring := fmt.Sprintf("graphql: panic occurred: %v", value)
	log.Println(errorstring)
}

type PanicHandler struct{}

func (p PanicHandler) MakePanicError(ctx context.Context, value interface{}) *gqlerrors.QueryError {
	return &gqlerrors.QueryError{
		Message: fmt.Sprintf("something went wrong processing your request"),
	}
}
