package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
)

type UserResolver struct {
	c    *graphql_context.Context
	user *models.User
}

func NewUserResolver(c *graphql_context.Context, user *models.User) (*UserResolver, error) {
	if user == nil {
		return nil, nil
	}

	return &UserResolver{c: c, user: user}, nil
}

func (r *UserResolver) UID() (graphql.ID, error) {
	id := graphql.ID(r.user.UID)
	return id, nil
}

func (r *UserResolver) DisplayName() (*string, error) {
	return nil, nil
}

func (r *UserResolver) Email() (*string, error) {
	return nil, nil
}

func (r *UserResolver) PhotoURL() (*string, error) {
	return nil, nil
}

func (b *BaseQuery) Me(ctx context.Context) (*UserResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser()
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	return NewUserResolver(c, currentUser)
}

type UserArgs struct {
	Uid graphql.ID
}

func (b *BaseQuery) User(ctx context.Context, args *UserArgs) (*UserResolver, error) {
	// @todo: implement me
	return nil, nil
}
