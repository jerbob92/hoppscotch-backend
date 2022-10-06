package resolvers

import (
	"context"
	"errors"
	"sync"

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
	if u.user.DisplayName == "" {
		return nil, nil
	}
	return &u.user.DisplayName, nil
}

func (u *UserResolver) Email() (*string, error) {
	if u.user.Email == "" {
		return nil, nil
	}
	return &u.user.Email, nil
}

func (u *UserResolver) PhotoURL() (*string, error) {
	if u.user.PhotoURL == "" {
		return nil, nil
	}
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

type UserSubscription struct {
	Lock        sync.Mutex
	UserDeleted map[string]chan *UserResolver
}

type UserSubscriptions struct {
	Subscriptions map[uint]*UserSubscription
	Lock          sync.Mutex
}

func (t *UserSubscriptions) EnsureChannel(channel uint) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if _, ok := t.Subscriptions[channel]; !ok {
		t.Subscriptions[channel] = &UserSubscription{
			Lock:        sync.Mutex{},
			UserDeleted: map[string]chan *UserResolver{},
		}
	}
}

var userSubscriptions = UserSubscriptions{
	Subscriptions: map[uint]*UserSubscription{},
	Lock:          sync.Mutex{},
}

func (b *BaseQuery) DeleteUser(ctx context.Context) (bool, error) {
	c := b.GetReqC(ctx)
	user, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return false, err
	}

	db := c.GetDB()
	err = db.Delete(user).Error
	if err != nil {
		return false, err
	}

	resolver, err := NewUserResolver(c, user)
	if err != nil {
		return false, err
	}

	go func() {
		userSubscriptions.EnsureChannel(user.ID)

		userSubscriptions.Subscriptions[user.ID].Lock.Lock()
		defer userSubscriptions.Subscriptions[user.ID].Lock.Unlock()
		for i := range userSubscriptions.Subscriptions[user.ID].UserDeleted {
			userSubscriptions.Subscriptions[user.ID].UserDeleted[i] <- resolver
		}
	}()

	return true, nil
}

func (b *BaseQuery) UserDeleted(ctx context.Context) (<-chan *UserResolver, error) {
	c := b.GetReqC(ctx)

	user, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	userSubscriptions.EnsureChannel(user.ID)

	notificationChannel := make(chan *UserResolver)
	subID := RandString(32)
	userSubscriptions.Subscriptions[user.ID].Lock.Lock()
	defer userSubscriptions.Subscriptions[user.ID].Lock.Unlock()
	userSubscriptions.Subscriptions[user.ID].UserDeleted[subID] = notificationChannel

	go func() {
		select {
		case <-ctx.Done():
			userSubscriptions.Subscriptions[user.ID].Lock.Lock()
			defer userSubscriptions.Subscriptions[user.ID].Lock.Unlock()
			close(userSubscriptions.Subscriptions[user.ID].UserDeleted[subID])
			delete(userSubscriptions.Subscriptions[user.ID].UserDeleted, subID)
			return
		}
	}()

	return notificationChannel, nil
}
