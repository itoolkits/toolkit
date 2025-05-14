// common response

package ctxt

import (
	"net/http"
)

type Response struct {
	SCode     int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id"`
	Cost      int64  `json:"cost"`
}

// NewOKResponse create ok response
func NewOKResponse() *Response {
	return &Response{
		SCode: http.StatusOK,
		Msg:   "",
	}
}
