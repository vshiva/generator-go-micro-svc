<%=licenseText%>
package auth

import (
	"fmt"
	"net/http"
	"net/url"
)

const cookieName = "express.sid"

// TokenMiddleware unifies an auth token to always be in the Authorization header.
//
// Supported inputs are (in order of preference):
//   1. Authorization header          // Authorization: Bearer abcdef
//   2. token query string parameter  // /api/doSomething?token=abcdef
//   3. express.sid cookie            // Cookie:express.sid=abdef;
//
// The resulting auth header will be `Authorization: <token>`. With cookies the result
// will be `Authorization: Cookie <value>`. If none of the authentication methods
// are given the Authorization header is not set.
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Auth header already set, pass it along without modifications
			next.ServeHTTP(w, r)
			return
		}

		token := r.URL.Query().Get("token")
		if token != "" {
			// Token was set in query string /api/service?token=ABCDEF
			// Set auth header with Bearer
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			next.ServeHTTP(w, r)
			return
		}

		authCookie, err := r.Cookie(cookieName)
		if err == nil {
			// A cookie was passed in, insert it as an auth header instead
			u, err := url.QueryUnescape(authCookie.Value)
			if err == nil {
				// The `Cookie` prefix tells the other side that this value
				// came from a cookie and should be decoded differently
				r.Header.Set("Authorization", fmt.Sprintf("Cookie %s", u))
			}
		}
		next.ServeHTTP(w, r)
	})
}
