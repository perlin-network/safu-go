package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/perlin-network/safu-go/log"
	"github.com/pkg/errors"
)

const (
	MaxRequestBodySize = 4 * 1024 * 1024
)

var (
	// ErrMsgBodyNil occurs when request body is empty
	ErrMsgBodyNil = errors.New("message body is nil")
)

// requestContext represents a context for a request.
type requestContext struct {
	service  *service
	response http.ResponseWriter
	request  *http.Request
}

// ErrorResponse is a payload when there is an error
type ErrorResponse struct {
	StatusCode int         `json:"status"`
	Error      interface{} `json:"error,omitempty"`
}

// readJSON decodes a HTTP requests JSON body into a struct.
// Can call this once per request
func (c *requestContext) readJSON(out interface{}) error {
	r := io.LimitReader(c.request.Body, MaxRequestBodySize)
	defer c.request.Body.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "bad request body")
	}

	if len(data) == 0 {
		return ErrMsgBodyNil
	}

	if err = json.Unmarshal(data, out); err != nil {
		return errors.Wrap(err, "malformed json")
	}
	return nil
}

// WriteJSON will write a given status code & JSON to a response.
// Should call this once per request
func (c *requestContext) WriteJSON(status int, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		c.WriteJSON(http.StatusInternalServerError, "server error")
		return
	}
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status)
	c.response.Write(out)
}

// wrap applies middleware to a HTTP request handler.
func (s *service) wrap(handler func(*requestContext) (int, interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &requestContext{
			service:  s,
			response: w,
			request:  r,
		}

		// recover from panics
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("IP", r.RemoteAddr).
					Str("path", r.URL.EscapedPath()).
					Msgf("An error occured from the API: %s", string(debug.Stack()))

				// return a 500 on a panic
				c.WriteJSON(http.StatusInternalServerError, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err,
				})
			}
		}()

		// call the handler
		statusCode, data, err := handler(c)

		// write the result
		if err != nil {
			log.Warn().
				Interface("IP", r.RemoteAddr).
				Str("path", r.URL.EscapedPath()).
				Interface("statusCode", statusCode).
				Msgf("An error occured from the API: %+v", err)

			c.WriteJSON(statusCode, ErrorResponse{
				StatusCode: statusCode,
				Error:      err.Error(),
			})
		} else {
			log.Debug().
				Interface("IP", r.RemoteAddr).
				Str("path", r.URL.EscapedPath()).
				Msg(" ")

			c.WriteJSON(statusCode, data)
		}
	}
}
