package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/bitgaming/gin-router-stub/api"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func middleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		sessions.Sessions("session", sessions.NewCookieStore([]byte("12345678"))),
	}
}

func TestPostUser(t *testing.T) {
	data := `{"username":"tedwen"}`
	RunSimplePost("/login", "/login", data,
		func(c *gin.Context) {
			api.PostUser(c)
		},
		func(r *httptest.ResponseRecorder) {
			if r.Code != 201 {
				t.Error(r)
			}
		})
}

func TestGetUser(t *testing.T) {
	RunWithMiddlewares("/login", "GET", "/login", "", middleWares(),
		func(c *gin.Context) {
			//prepare session
			api.GetUser(c)
		},
		func(r *httptest.ResponseRecorder) {
			t.Logf("r=%+v", r)
		})
}
