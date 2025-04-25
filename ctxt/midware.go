// http route middleware

package ctxt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthFunc - auth func
func AuthFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := GetCTX(ctx)
		ak := ctx.GetHeader("AK")
		auth := ctx.GetHeader("Authorization")
		tsStr := ctx.GetHeader("Timestamp")
		if tsStr == "" {
			reqCtx.SetErrorResponse(http.StatusBadRequest, http.StatusBadRequest, fmt.Sprintf("timestamp can not blank"))
			ctx.Abort()
			return
		}
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			reqCtx.SetErrorResponse(http.StatusBadRequest, http.StatusBadRequest, fmt.Sprintf("timestamp format error,  %v", err))
			ctx.Abort()
			return
		}

		err = MD5Check(ak, auth, ts)
		if err != nil {
			reqCtx.SetErrorResponse(http.StatusUnauthorized, http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// ErrorFunc - error func
func ErrorFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				recoverPanic(c, err)
				c.Abort()
				return
			}
		}(c)
		c.Next()
	}
}

// ResponseFunc - response function
func ResponseFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		reqCtx := GetCTX(c)
		if reqCtx.HTTPCode == 0 {
			reqCtx.HTTPCode = http.StatusOK
		}

		if reqCtx.Response == nil {
			reqCtx.Response = NewOKResponse()
		}

		if reqCtx.Response.SCode < 600 {
			reqCtx.HTTPCode = reqCtx.Response.SCode
		}

		reqCtx.Response.Cost = time.Now().UnixMilli() - start.UnixMilli()

		c.JSON(reqCtx.HTTPCode, reqCtx.Response)
	}
}

// recoverPanic - recover panic error
func recoverPanic(c *gin.Context, err any) {
	reqCtx := GetCTX(c)

	reqCtx.TraceStack = fmt.Sprintf("[Error Stack]  [%s]  [%v] \n%s\n", reqCtx.RequestID, err, gatherStackInfo())

	slog.Error("find unexpect error, when handle the request", "requestID", reqCtx.RequestID, "error", err)

	msg := fmt.Sprintf("unexpect error, %v", err)

	reqCtx.SetErrorResponse(500, 500, msg)
}

// gatherStackInfo - Get call stack when error
func gatherStackInfo() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return strings.ReplaceAll(string(buf), "\n", " | ")
}

const (
	RequestBodyParam = "RequestBodyParam"
	AcLogMaxLen      = 2048
)

type AccessLogVO struct {
	RequestID     string    `json:"request_id"`
	ServerHost    string    `json:"server_host"`
	ClientHost    string    `json:"client_host"`
	StartTime     time.Time `json:"time"`
	HTTPMethod    string    `json:"http_method"`
	HTTPURL       string    `json:"http_url"`
	HTTPHeader    string    `json:"http_header"`
	HTTPParam     string    `json:"http_param"`
	HTTPCode      int       `json:"http_code"`
	ResponseCode  int       `json:"response_code"`
	ResponseParam string    `json:"response_param"`
	ErrorMsg      string    `json:"error_msg"`
	ErrorTrace    string    `json:"error_trace"`
	DurationMS    int64     `json:"duration_ms"`
}

// LogAttr - log attr
func (a *AccessLogVO) LogAttr() []any {
	rst := []any{}
	rst = append(rst, "request_id", a.RequestID)
	rst = append(rst, "server_host", a.ServerHost)
	rst = append(rst, "client_host", a.ClientHost)
	rst = append(rst, "start_time", a.StartTime.Format(time.DateTime))
	rst = append(rst, "http_method", a.HTTPMethod)
	rst = append(rst, "http_url", a.HTTPURL)
	rst = append(rst, "http_header", a.HTTPHeader)
	rst = append(rst, "http_param", a.HTTPParam)
	rst = append(rst, "http_code", a.HTTPCode)
	rst = append(rst, "response_code", a.ResponseCode)
	rst = append(rst, "response_param", a.ResponseParam)
	rst = append(rst, "error_msg", a.ErrorMsg)
	rst = append(rst, "error_trace", a.ErrorTrace)
	rst = append(rst, "duration_ms", a.DurationMS)
	return rst
}

// ACLog - access log
func ACLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// API context
		reqCtx := GetCTX(c)

		// Access Log
		accessLog := &AccessLogVO{}

		startTime := time.Now()

		// Body bytes
		bts, err := io.ReadAll(c.Request.Body)
		var cbts []byte
		if len(bts) > 0 && err == nil {
			cbts = make([]byte, len(bts))
			copy(cbts, bts)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bts))
		}
		requestParams := getParams(c, cbts)
		if len(requestParams) > 0 {
			requestParamStr, err := json.Marshal(requestParams)
			if err == nil && len(requestParamStr) > 0 {
				accessLog.HTTPParam = string(requestParamStr)
			}
		}

		accessLog.RequestID = reqCtx.RequestID
		accessLog.ClientHost = getClientIP(c)
		accessLog.ServerHost = getHostName()

		accessLog.HTTPMethod = c.Request.Method
		accessLog.HTTPURL = c.Request.URL.Path
		accessLog.StartTime = startTime

		accessLog.HTTPHeader = getHeaderParam(c)

		// Process request
		c.Next()

		// Collect log info
		cost := time.Since(startTime)

		accessLog.DurationMS = int64(cost) / 1e6

		rsps := reqCtx.Response

		accessLog.ErrorMsg = rsps.Msg
		accessLog.ResponseCode = rsps.SCode
		accessLog.HTTPCode = reqCtx.HTTPCode

		accessLog.ErrorTrace = reqCtx.TraceStack

		retData := rsps.Data
		if retData != nil {
			retBts, _ := json.Marshal(retData)
			resBitsLen := len(retBts)
			if resBitsLen > AcLogMaxLen {
				resBitsLen = AcLogMaxLen
			}
			accessLog.ResponseParam = string(retBts[:resBitsLen])
		}

		switch {
		case accessLog.HTTPCode >= 400 && accessLog.HTTPCode <= 499:
			{
				AccessLogger.Warn("request error", accessLog.LogAttr()...)
			}
		case accessLog.HTTPCode >= 500:
			{
				AccessLogger.Error("server error", accessLog.LogAttr()...)
			}
		default:
			{
				AccessLogger.Info("request success", accessLog.LogAttr()...)
			}
		}
	}
}

// getParams - get params from context, support query/params/body
func getParams(c *gin.Context, body []byte) map[string]string {
	params := make(map[string]string)
	// query data
	values := c.Request.URL.Query()
	if len(values) > 0 {
		for k, vs := range values {
			if len(vs) > 0 {
				params[k] = vs[0]
			}
		}
	}

	// param data
	if len(c.Params) > 0 {
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}
	}

	// body data
	if len(body) > 0 {
		params[RequestBodyParam] = string(body)
	}

	return params
}

// getHeaderParam - get header parma as json string
func getHeaderParam(c *gin.Context) string {
	if len(c.Request.Header) > 0 {
		hj, err := json.Marshal(c.Request.Header)
		if err == nil && len(hj) > 0 {
			return string(hj)
		}
	}
	return ""
}

// getClientIP - get client ip
func getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()
	if clientIP != "" {
		return clientIP
	}

	// No client ip header or invalid header, return 0.0.0.0
	return "0.0.0.0"
}

// getHostName - get machine host name
func getHostName() string {
	hostName, have := os.LookupEnv("HOSTNAME")
	if !have || hostName == "" {
		return "localhost"
	}
	return hostName
}
