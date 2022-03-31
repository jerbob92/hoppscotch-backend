package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strconv"
)

type ShortcodeResolver struct {
	c         *graphql_context.Context
	shortcode *models.Shortcode
}

func NewShortcodeResolver(c *graphql_context.Context, shortcode *models.Shortcode) (*ShortcodeResolver, error) {
	if shortcode == nil {
		return nil, nil
	}

	return &ShortcodeResolver{c: c, shortcode: shortcode}, nil
}

func (r *ShortcodeResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.FormatInt(r.shortcode.ID, 10))
	return id, nil
}

func (r *ShortcodeResolver) Request() (string, error) {
	return r.shortcode.Request, nil
}

type ShortcodeArgs struct {
	Code graphql.ID
}

func (b *BaseQuery) Shortcode(ctx context.Context, args *ShortcodeArgs) (*ShortcodeResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateShortcodeArgs struct {
	Request string
}

func (b *BaseQuery) CreateShortcode(ctx context.Context, args *CreateShortcodeArgs) (*ShortcodeResolver, error) {
	// @todo: implement me
	return nil, nil
}
