// context struct

package ctxt

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestContext = "REQUEST_CONTEXT"
)

type RqstCtx struct {
	RequestID string

	Response *Response

	HTTPCode int

	TraceStack string
}

// initRequestCTX - init request ctx
func (r *RqstCtx) init(c *gin.Context) {
	r.setRequestID(c)
	c.Set(requestContext, r)
}

// setRequestID - set request id
func (r *RqstCtx) setRequestID(c *gin.Context) {
	requestID := c.GetHeader("requestID")
	if requestID != "" {
		r.RequestID = requestID
		return
	}
	r.RequestID = uuid.New().String()
	c.Header("requestID", r.RequestID)
}

// SetHTTPCode - set http code
func (r *RqstCtx) SetHTTPCode(code int) {
	r.HTTPCode = code
}

// SetErrorResponse - set error response
func (r *RqstCtx) SetErrorResponse(httpCode int, sCode int, msg string) {
	r.SetHTTPCode(httpCode)
	r.Response = &Response{
		SCode:     sCode,
		Msg:       msg,
		RequestID: r.RequestID,
	}
}

// SetStatusResponse - set status response
func (r *RqstCtx) SetStatusResponse(sCode int, msg string, args ...any) {
	r.SetHTTPCode(http.StatusOK)
	msg = fmt.Sprintf(msg, args...)
	r.Response = &Response{
		SCode:     sCode,
		Msg:       msg,
		RequestID: r.RequestID,
	}
}

// SetOKResponse - set ok response
func (r *RqstCtx) SetOKResponse() {
	r.Response = &Response{
		SCode:     http.StatusOK,
		Msg:       "",
		RequestID: r.RequestID,
	}
}

// SetData - set data
func (r *RqstCtx) SetData(data any) {
	r.Response = &Response{
		SCode:     http.StatusOK,
		Msg:       "",
		Data:      data,
		RequestID: r.RequestID,
	}
}

// GetCTX - get context
func GetCTX(c *gin.Context) *RqstCtx {
	reqCTXObj, have := c.Get(requestContext)
	if !have {
		reqCTX := &RqstCtx{}
		reqCTX.init(c)
		return reqCTX
	}
	return reqCTXObj.(*RqstCtx)
}
