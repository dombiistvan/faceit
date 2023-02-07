package helpers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

type RequestHandler struct {
	activeRequests map[string]bool
	sm             *sync.RWMutex
}

var rh = &RequestHandler{
	activeRequests: make(map[string]bool),
	sm:             new(sync.RWMutex),
}

// json middleware for create/update requests
func JsonMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Request().Header.Get("Content-type"), "application/json") {
			return echo.NewHTTPError(http.StatusBadRequest, "expected request content-type is application/json")
		}

		return next(c)
	}
}

// too many requests middleware preveting server/app from spamming
func TooManyRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var headers []string = []string{}
		var headerValues []string = []string{}
		body, err := io.ReadAll(c.Request().Body)
		ErrChan <- err
		// Replace the body with a new reader after reading from the original
		c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

		for k, _ := range c.Request().Header {
			headers = append(headers, k)
		}

		sort.Strings(headers)
		for _, k := range headers {
			headerValues = append(headerValues, fmt.Sprintf("header:%s=%v", k, c.Request().Header.Get(k)))
		}

		encodedForm := c.Request().Form.Encode()

		requestHash := sha256.New()
		requestHash.Write([]byte(strings.Join(headers, "&")))
		requestHash.Write([]byte(encodedForm))
		requestHash.Write(body)

		requestKey := fmt.Sprintf("%x", requestHash.Sum(nil))

		rh.sm.Lock()
		_, ok := rh.activeRequests[requestKey]
		if ok {
			return c.NoContent(http.StatusTooManyRequests)
		}
		rh.activeRequests[requestKey] = true
		rh.sm.Unlock()

		defer func(rk string) {
			rh.sm.Lock()
			delete(rh.activeRequests, requestKey)
			rh.sm.Unlock()
		}(requestKey)

		return next(c)
	}
}
