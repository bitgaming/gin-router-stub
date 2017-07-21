package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//RequestFunc is request handling func type to replace gin.HandlerFunc
type RequestFunc func(*gin.Context)

//ResponseFunc is response handling func type
type ResponseFunc func(*httptest.ResponseRecorder)

type RequestConfig struct {
	Method      string
	RoutePath   string
	Path        string
	Body        string
	Headers     map[string]string
	Middlewares []gin.HandlerFunc
	Handler     RequestFunc
	Finaliser   ResponseFunc
}

func injectClock(c *gin.Context) {
	c.Set("freezedCurrentTime", time.Now().UTC())
	c.Next()
}

func RunRequest(rc RequestConfig) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(injectClock)

	if rc.Middlewares != nil && len(rc.Middlewares) > 0 {
		for _, mw := range rc.Middlewares {
			r.Use(mw)
		}
	}

	qs := ""
	if strings.Contains(rc.Path, "?") {
		ss := strings.Split(rc.Path, "?")
		rc.Path = ss[0]
		qs = ss[1]
	}

	body := bytes.NewBufferString(rc.Body)

	req, _ := http.NewRequest(rc.Method, rc.Path, body)

	if len(qs) > 0 {
		req.URL.RawQuery = qs
	}

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	} else if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	r.Handle(rc.Method, rc.RoutePath, func(c *gin.Context) {
		//change argument if necessary here
		rc.Handler(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if rc.Finaliser != nil {
		rc.Finaliser(w)
	}
}

func RunSimpleGet(routePath string, path string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "GET",
		RoutePath: routePath,
		Path:      path,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func RunSimplePost(routePath string, path, body string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "POST",
		RoutePath: routePath,
		Path:      path,
		Body:      body,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func RunGetWithHeaders(routePath string, path string, headers map[string]string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "GET",
		RoutePath: routePath,
		Path:      path,
		Headers:   headers,
		Handler:   handler,
		Finaliser: reply,
	}
	RunRequest(rc)
}

func RunPostWithHeaders(routePath string, path string, body string, headers map[string]string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:    "POST",
		Body:      body,
		RoutePath: routePath,
		Path:      path,
		Headers:   headers,
		Handler:   handler,
		Finaliser: reply,
	}

	RunRequest(rc)
}

func RunGetWithMiddlewares(routePath string, path string, mws []gin.HandlerFunc, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      "GET",
		RoutePath:   routePath,
		Path:        path,
		Middlewares: mws,
		Handler:     handler,
		Finaliser:   reply,
	}
	RunRequest(rc)
}

func RunPostWithMiddlewares(routePath string, path, body string, mws []gin.HandlerFunc, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      "POST",
		RoutePath:   routePath,
		Path:        path,
		Body:        body,
		Middlewares: mws,
		Handler:     handler,
		Finaliser:   reply,
	}
	RunRequest(rc)
}

func RunWithMiddlewares(routePath string, method, path, body string, mws []gin.HandlerFunc, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      method,
		RoutePath:   routePath,
		Path:        path,
		Middlewares: mws,
		Handler:     handler,
		Finaliser:   reply,
	}
	if method == "POST" || method == "PUT" {
		rc.Body = body
	}
	RunRequest(rc)
}

func RunWithHeaderMiddlewares(routePath string, method, path, body string, mws []gin.HandlerFunc, headers map[string]string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      method,
		RoutePath:   routePath,
		Path:        path,
		Headers:     headers,
		Middlewares: mws,
		Handler:     handler,
		Finaliser:   reply,
	}
	if method == "POST" || method == "PUT" {
		rc.Body = body
	}
	RunRequest(rc)
}

func VerifyResponse(r *httptest.ResponseRecorder, code int, data map[string]interface{}) (bool, error) {
	if r.Code != code {
		s := fmt.Sprintf("Code returned:%d != expected:%d", r.Code, code)
		return false, errors.New(s)
	}
	var rd map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return false, err
	}
	if !reflect.DeepEqual(rd, data) {
		s := fmt.Sprintf("R.Body:%+v", rd)
		return false, errors.New(s)
	}
	return true, nil
}
