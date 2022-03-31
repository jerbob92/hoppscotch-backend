package resolvers

import (
	"context"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
)

// BaseQuery is the query resolvers
type BaseQuery struct {
	c *graphql_context.Context // Only set/overwrite this context when constructing new baseQueries.
}

// GetReqC returns the request context
func (b *BaseQuery) GetReqC(ctx context.Context) *graphql_context.Context {
	if b.c == nil {
		interf := ctx.Value("graphqlC")
		c := interf.(*graphql_context.Context)
		return c.Clone()
	}

	// Clone the request context so mutationCompany and queryCompany only have effects on the inner request of mutationCompany and queryCompany
	return b.c.Clone()
}
