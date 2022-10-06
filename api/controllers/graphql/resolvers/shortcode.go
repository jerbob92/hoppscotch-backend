package resolvers

import (
	"context"
	"crypto/rand"
	"errors"
	"strconv"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
	"github.com/sanae10001/graphql-go-extension-scalars"
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

func (r *ShortcodeResolver) CreatedOn() (scalars.DateTime, error) {
	return *scalars.NewDateTime(r.shortcode.CreatedAt), nil
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

	resolver, err := NewShortcodeResolver(c, newShortCode)
	if err != nil {
		return nil, err
	}

	bus.Publish("user:"+strconv.Itoa(int(currentUser.ID))+":shortcodes:created", resolver)

	return resolver, nil
}

type RevokeShortcodeArgs struct {
	Code graphql.ID
}

func (b *BaseQuery) RevokeShortcode(ctx context.Context, args *RevokeShortcodeArgs) (bool, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return false, err
	}

	shortcode := &models.Shortcode{}
	db := c.GetDB()
	err = db.Model(&models.Shortcode{}).Where("code = ?", args.Code).First(shortcode).Error
	if err != nil {
		return false, err
	}

	if shortcode.UserID != currentUser.ID {
		return false, errors.New("you do not have access to this shortcode")
	}

	err = db.Delete(shortcode).Error
	if err != nil {
		return false, err
	}

	resolver, err := NewShortcodeResolver(c, shortcode)
	if err != nil {
		return false, err
	}

	bus.Publish("user:"+strconv.Itoa(int(currentUser.ID))+":shortcodes:revoked", resolver)

	return true, nil
}

type MyShortcodeArgs struct {
	Cursor *graphql.ID
}

func (b BaseQuery) MyShortcodes(ctx context.Context, args *MyShortcodeArgs) ([]*ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	shortcodes := []*models.Shortcode{}
	db := c.GetDB()
	query := db.Model(&models.Shortcode{}).Where("user_id = ?", currentUser.ID)
	if args.Cursor != nil && *args.Cursor != "" {
		query.Where("id > ?", args.Cursor)
	}
	err = query.Find(&shortcodes).Error
	if err != nil {
		return nil, err
	}

	shortcodesResolvers := []*ShortcodeResolver{}
	for i := range shortcodes {
		newResolver, err := NewShortcodeResolver(c, shortcodes[i])
		if err != nil {
			return nil, err
		}
		shortcodesResolvers = append(shortcodesResolvers, newResolver)
	}

	return shortcodesResolvers, nil
}

func (b *BaseQuery) MyShortcodesCreated(ctx context.Context) (<-chan *ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	notificationChannel := make(chan *ShortcodeResolver)
	eventHandler := func(resolver *ShortcodeResolver) {
		notificationChannel <- resolver
	}
	err = subscribeUntilDone(ctx, "user:"+strconv.Itoa(int(currentUser.ID))+":shortcodes:created", eventHandler)
	if err != nil {
		return nil, err
	}

	return notificationChannel, nil
}

func (b *BaseQuery) MyShortcodesRevoked(ctx context.Context) (<-chan *ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	notificationChannel := make(chan *ShortcodeResolver)
	eventHandler := func(resolver *ShortcodeResolver) {
		notificationChannel <- resolver
	}
	err = subscribeUntilDone(ctx, "user:"+strconv.Itoa(int(currentUser.ID))+":shortcodes:revoked", eventHandler)
	if err != nil {
		return nil, err
	}

	return notificationChannel, nil
}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
