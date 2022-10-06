package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/resolvers"
	"github.com/jerbob92/hoppscotch-backend/schema"

	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-go/trace"
	_ "github.com/sanae10001/graphql-go-extension-scalars"
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
		graphql.Tracer(graphqlTracer{}),
	)

	Handler = &relay.Handler{Schema: parsedSchema}
}

type graphqlTracer struct{}

func (t graphqlTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, varTypes map[string]*introspection.Type) (context.Context, trace.TraceQueryFinishFunc) {
	//log.Println("Trace query")
	//log.Println(queryString)
	//log.Println(operationName)
	//log.Println(variables)
	return ctx, func(errs []*gqlerrors.QueryError) {}
}

func (t graphqlTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, trace.TraceFieldFinishFunc) {
	//log.Println("Trace field")
	//log.Println(label)
	//log.Println(typeName)
	//log.Println(fieldName)
	//log.Println(args)
	return ctx, func(err *gqlerrors.QueryError) {}
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
