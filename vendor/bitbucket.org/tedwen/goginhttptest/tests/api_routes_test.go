package tests

import (
	"testing"
	"net/http/httptest"
	"bitbucket.org/tedwen/goginhttptest/api"
)

func TestPostUser(t *testing.T) {
	data := `{"username":"tedwen"}`
	RunSimplePost("/login", data,
		func(c *gin.Context){
			api.PostUser(c)
		},
		func(r *httptest.ResponseRecorder) {
			if r.Code != 201 {
				t.Error(r)
			}
		})
}

func TestGetUser(t *testing.T) {
	RunWithMiddlewares("GET", "/login", "",
		func(c *gin.Context){
			//prepare session
			api.GetUser(c)
		},
		func(r *httptest.ResponseRecorder){
			t.Logf("r=%+v", r)
		})
}
