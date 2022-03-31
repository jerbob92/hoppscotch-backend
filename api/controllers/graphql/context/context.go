package context

import (
	"errors"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func (c *Context) GetUser() (*models.User, error) {
	if c.ReqUser != nil {
		return c.ReqUser, nil
	}

	if c.GinContext == nil {
		return nil, errors.New("could not load user from request")
	}

	header := c.GinContext.Request.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		header = header[7:]
	}

	// @todo: try to get user.

	return nil, nil
}
