package resolvers

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"gorm.io/gorm"
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
	id := graphql.ID(r.shortcode.Code)
	return id, nil
}

func (r *ShortcodeResolver) Request() (string, error) {
	return r.shortcode.Request, nil
}

type ShortcodeArgs struct {
	Code graphql.ID
}

func (b *BaseQuery) Shortcode(ctx context.Context, args *ShortcodeArgs) (*ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()
	shortcode := &models.Shortcode{}
	err := db.Model(&models.Shortcode{}).Where("code = ?", args.Code).First(shortcode).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this shortcode")
	}
	if err != nil {
		return nil, err
	}

	return NewShortcodeResolver(c, shortcode)
}

type CreateShortcodeArgs struct {
	Request string
}

func (b *BaseQuery) CreateShortcode(ctx context.Context, args *CreateShortcodeArgs) (*ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	db := c.GetDB()
	newShortCode := &models.Shortcode{
		Code:    RandString(12),
		Request: args.Request,
		UserID:  currentUser.ID,
	}

	err = db.Save(newShortCode).Error
	if err != nil {
		return nil, err
	}

	return NewShortcodeResolver(c, newShortCode)
}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const alphanumlower = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func RandStringLower(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanumlower[b%byte(len(alphanumlower))]
	}
	return string(bytes)
}
