package context

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jerbob92/hoppscotch-backend/fb"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Context is a global request context
type Context struct {
	ReqUser     *models.User
	GinContext  *gin.Context
	loggingMeta map[string]interface{}
	locking     sync.Mutex

	// DisableResponses is mainly used for the graphql routes because the library handles error messages and we don't want to return custom errors
	DisableResponses bool
}

func GetContext(c *gin.Context) *Context {
	return &Context{
		GinContext:  c,
		loggingMeta: map[string]interface{}{},
		locking:     sync.Mutex{},
	}
}

func (c *Context) SetLoggingMetaValue(key string, value interface{}) {
	c.locking.Lock()
	defer c.locking.Unlock()
	c.loggingMeta[key] = value
}

func (c *Context) Clone() *Context {
	c.locking.Lock()
	defer c.locking.Unlock()

	newLoggingMeta := map[string]interface{}{}
	for k, v := range c.loggingMeta {
		newLoggingMeta[k] = v
	}

	newContext := &Context{
		GinContext:       c.GinContext,
		loggingMeta:      newLoggingMeta,
		locking:          sync.Mutex{},
		DisableResponses: c.DisableResponses,
	}

	if c.ReqUser != nil {
		newContext.ReqUser = &*c.ReqUser
	}

	return newContext
}

func (c *Context) LogErr(err error, meta ...map[string]interface{}) {
	data := log.Fields{
		"type": "USER_API_ERROR",
	}

	c.locking.Lock()
	for key, value := range c.loggingMeta {
		data[key] = value
	}
	c.locking.Unlock()

	for _, metaItem := range meta {
		for key, value := range metaItem {
			data[key] = value
		}
	}
	log.WithFields(data).Error(err)
}

func (c *Context) GetDB() *gorm.DB {
	if c.GinContext == nil {
		return nil
	}
	db, exists := c.GinContext.Get("DB")
	if !exists {
		return nil
	}

	return db.(*gorm.DB)
}

func (c *Context) GetUser(ctx context.Context) (*models.User, error) {
	if c.ReqUser != nil {
		return c.ReqUser, nil
	}

	if c.GinContext == nil {
		return nil, errors.New("could not load user from request")
	}

	header := c.GinContext.Request.Header.Get("Authorization")

	// Fallback for subscription header.
	if header == "" {
		headerJSON, ok := ctx.Value("Header").(json.RawMessage)
		if ok {
			type initMessagePayload struct {
				Authorization string `json:"authorization"`
			}

			var initMsg initMessagePayload
			if err := json.Unmarshal(headerJSON, &initMsg); err != nil {
				return nil, err
			}

			header = initMsg.Authorization
		}
	}
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		header = header[7:]
	}

	token, err := fb.FBAuth.VerifyIDToken(context.Background(), header)
	if err != nil {
		return nil, fmt.Errorf("could not validate ID token: %s", header)
	}

	db := c.GetDB()
	if db == nil {
		return nil, errors.New("can't get DB")
	}

	existingUser := &models.User{}
	err = db.Where("fb_uid = ?", token.UID).First(existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		newUser := &models.User{
			FBUID:       token.UID,
			DisplayName: "",
			PhotoURL:    "",
		}

		email, ok := token.Claims["email"].(string)
		if ok {
			newUser.Email = email
		}

		userObj, err := fb.FBAuth.GetUser(ctx, token.UID)
		if err != nil {
			return nil, err
		}

		newUser.DisplayName = userObj.UserInfo.DisplayName
		newUser.PhotoURL = userObj.UserInfo.PhotoURL

		err = db.Create(newUser).Error
		if err != nil {
			return nil, err
		}

		c.ReqUser = newUser

		return newUser, nil
	}

	if err != nil {
		return nil, err
	}

	c.ReqUser = existingUser

	return existingUser, nil
}
