package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

// request handling func type to replace gin.HandlerFunc
type RequestFunc func(*gin.Context)

// response handling func type
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

func RunRequest(rc RequestConfig) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

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

func RunWithMiddlewares(routePath string, method, path, body string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      method,
		RoutePath:   routePath,
		Path:        path,
		Middlewares: MiddleWares(),
		Handler:     handler,
		Finaliser:   reply,
	}
	if method == "POST" || method == "PUT" {
		rc.Body = body
	}
	RunRequest(rc)
}

func RunWithHeaderMiddlewares(routePath string, method, path, body string, headers map[string]string, handler RequestFunc, reply ResponseFunc) {
	rc := RequestConfig{
		Method:      method,
		RoutePath:   routePath,
		Path:        path,
		Headers:     headers,
		Middlewares: MiddleWares(),
		Handler:     handler,
		Finaliser:   reply,
	}
	if method == "POST" || method == "PUT" {
		rc.Body = body
	}
	RunRequest(rc)
}

func MiddleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		sessions.Sessions("session", sessions.NewCookieStore([]byte("12345678"))),
	}
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
