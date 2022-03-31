package resolvers

import (
	"context"
	"errors"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
	"gorm.io/gorm"
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

func (u *UserResolver) UID() (graphql.ID, error) {
	id := graphql.ID(u.user.FBUID)
	return id, nil
}

func (u *UserResolver) DisplayName() (*string, error) {
	return &u.user.DisplayName, nil
}

func (u *UserResolver) Email() (*string, error) {
	return &u.user.Email, nil
}

func (u *UserResolver) PhotoURL() (*string, error) {
	return &u.user.PhotoURL, nil
}

func (b *BaseQuery) Me(ctx context.Context) (*UserResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
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
	c := b.GetReqC(ctx)
	_, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()
	existingUser := &models.User{}
	err = db.Where("fb_uid = ?", args.Uid).First(existingUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	}

	return NewUserResolver(c, existingUser)
}
