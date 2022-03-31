package responses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

// The serializable Error structure.
type Error struct {
	Result       string `json:"Result" description:"Whether the call was successful." enum:"success,error"`
	ErrorCode    int    `json:"Code" description:"The code of the error."`
	ErrorMessage string `json:"Message" description:"The human message of the error."`
	RequestId    string `json:"RequestId" description:"A request ID for support and debugging purposes."`
}

type RequestError struct {
	Code    int
	Message string
}

func (r RequestError) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

type RequestInternalError struct {
	OriginalError error
	Code          int
	Message       string
}

func (r RequestInternalError) Error() string {
	return fmt.Sprintf("%d: %s: %v", r.Code, r.Message, r.OriginalError)
}

type Success struct {
	Result    string      `json:"Result" description:"Whether the call was successful." enum:"success,error"`
	Data      interface{} `json:"Data"`
	RequestId string      `json:"RequestId" description:"A request ID for support and debugging purposes."`
}

func GraphQLInternalError(originalError error, requestID string) error {
	// Information for in the logging
	log.Error(fmt.Sprintf("error occurred, message: %v", originalError))

	// Information the user sees
	return fmt.Errorf("something went wrong processing your request: %s", requestID)
}

func GraphQLError(message string, requestID string) error {
	return fmt.Errorf("%s: %s", message, requestID)
}

func JSONAbort(c *gin.Context, code int, obj interface{}) {
	c.Abort()
	JSON(c, code, obj)
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func JSON(c *gin.Context, code int, obj interface{}) {
	if code >= 400 {
		returnCode, returnObject := RenderInternalError(c, code, obj)
		c.Render(returnCode, render.JSON{Data: returnObject})
		return
	}

	c.Render(code, render.JSON{Data: Success{
		Result: "success",
		Data:   obj,
	}})
}

func RenderInternalError(c *gin.Context, code int, obj interface{}) (int, interface{}) {
	errorObj := Error{
		Result:       "error",
		ErrorCode:    InputValidationError,
		ErrorMessage: "Request failed",
	}

	if requestError, ok := obj.(RequestError); ok {
		if requestError.Code != 0 {
			errorObj.ErrorCode = requestError.Code
		}
		if requestError.Message != "" {
			errorObj.ErrorMessage = requestError.Message
		}

		if errorObj.ErrorCode == InputValidationError {
			code = http.StatusBadRequest
		}

		if errorObj.ErrorCode == Unauthorized {
			code = http.StatusUnauthorized
		}
	}

	if requestError, ok := obj.(RequestInternalError); ok {
		errorObj.ErrorCode = InternalError
		if requestError.Code != 0 {
			errorObj.ErrorCode = requestError.Code
		}
		if requestError.Message != "" {
			errorObj.ErrorMessage = requestError.Message
		}
		code = http.StatusInternalServerError
	}

	return code, errorObj
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	errorObj := Error{
		Result:       "error",
		ErrorCode:    InternalError,
		ErrorMessage: "Request failed",
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, errorObj)
}
