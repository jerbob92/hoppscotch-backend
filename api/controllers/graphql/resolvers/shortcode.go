package resolvers

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"

	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"gorm.io/gorm"
)

type ShortcodeResolver struct {
	c         *graphql_context.Context
	shortcode *models.Shortcode
}

type ShortcodeSubscriptions struct {
	Lock                sync.Mutex
	MyShortcodesCreated map[string]chan *ShortcodeResolver
	MyShortcodesRevoked map[string]chan *ShortcodeResolver
}

type UserShortcodeSubscriptions struct {
	Subscriptions map[uint]*ShortcodeSubscriptions
	Lock          sync.Mutex
}

func (t *UserShortcodeSubscriptions) EnsureChannel(channel uint) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if _, ok := t.Subscriptions[channel]; !ok {
		t.Subscriptions[channel] = &ShortcodeSubscriptions{
			Lock:                sync.Mutex{},
			MyShortcodesCreated: map[string]chan *ShortcodeResolver{},
			MyShortcodesRevoked: map[string]chan *ShortcodeResolver{},
		}
	}
}

var userShortcodeSubscriptions = UserShortcodeSubscriptions{
	Subscriptions: map[uint]*ShortcodeSubscriptions{},
	Lock:          sync.Mutex{},
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

func (r *ShortcodeResolver) CreatedOn() (string, error) {
	return r.shortcode.CreatedAt.String(), nil
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

	go func() {
		userShortcodeSubscriptions.EnsureChannel(currentUser.ID)

		userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
		defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
		for i := range userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated {
			userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated[i] <- resolver
		}
	}()

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

	go func() {
		userShortcodeSubscriptions.EnsureChannel(currentUser.ID)

		userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
		defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
		for i := range userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated {
			userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated[i] <- resolver
		}
	}()

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

	userShortcodeSubscriptions.EnsureChannel(currentUser.ID)

	notificationChannel := make(chan *ShortcodeResolver)
	subID := RandString(32)
	userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
	defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
	userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated[subID] = notificationChannel

	go func() {
		select {
		case <-ctx.Done():
			userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
			defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
			close(userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated[subID])
			delete(userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesCreated, subID)
			return
		}
	}()

	return notificationChannel, nil
}

func (b *BaseQuery) MyShortcodesRevoked(ctx context.Context) (<-chan *ShortcodeResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	userShortcodeSubscriptions.EnsureChannel(currentUser.ID)

	notificationChannel := make(chan *ShortcodeResolver)
	subID := RandString(32)
	userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
	defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
	userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesRevoked[subID] = notificationChannel

	go func() {
		select {
		case <-ctx.Done():
			userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Lock()
			defer userShortcodeSubscriptions.Subscriptions[currentUser.ID].Lock.Unlock()
			close(userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesRevoked[subID])
			delete(userShortcodeSubscriptions.Subscriptions[currentUser.ID].MyShortcodesRevoked, subID)
			return
		}
	}()

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
