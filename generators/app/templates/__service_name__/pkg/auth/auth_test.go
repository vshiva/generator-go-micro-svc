<%=licenseText%>
package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mirror = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Just mirror back the auth header so we can look at the result
	w.Header().Set("Authorization", r.Header.Get("Authorization"))
})

func TestMiddleware_Auth_Blank(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	AuthTokenMiddleware(mirror).ServeHTTP(rr, req)

	assert.Equal(t, "", rr.Header().Get("Authorization"))
}

func TestMiddleware_Auth_Existing(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?token=query", nil)
	// auth header is set so it should be used
	req.Header.Set("Authorization", "auth")
	req.Header.Set("Cookie", "express.sid=value")
	AuthTokenMiddleware(mirror).ServeHTTP(rr, req)

	assert.Equal(t, "auth", rr.Header().Get("Authorization"))
}

func TestMiddleware_Auth_Query(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?token=token", nil)
	AuthTokenMiddleware(mirror).ServeHTTP(rr, req)

	assert.Equal(t, "Bearer token", rr.Header().Get("Authorization"))
}

func TestMiddleware_Auth_Cookie(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "express.sid=value")
	AuthTokenMiddleware(mirror).ServeHTTP(rr, req)

	assert.Equal(t, "Cookie value", rr.Header().Get("Authorization"))
}

func TestMiddleware_Auth_OtherCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "othercookie=value")
	AuthTokenMiddleware(mirror).ServeHTTP(rr, req)

	assert.Equal(t, "", rr.Header().Get("Authorization"))
}
