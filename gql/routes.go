package gql

import (
	"fmt"
	"net/http"

	"github.com/tastycrayon/blog-backend/middleware"
)

func RefreshTokenRoute(w http.ResponseWriter, r *http.Request) {
	response := `{ "ok": false, "access_token": "" }`
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("access-Control-Allow-Credentials", "true")

	refresh_token, err := r.Cookie("refresh_token")
	if refresh_token == nil || err != nil {
		w.WriteHeader(401)
		w.Write([]byte(response))
		return
	}
	userId, err := middleware.ParseToken(refresh_token.Value)
	// check database // check user._v
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(response))
		return
	}

	access_token, _, err := middleware.GenerateAccessToken(*userId)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(response))
		return
	}
	response = fmt.Sprintf(`{ "ok": true, "access_token": "%v" }`, *access_token)
	w.Write([]byte(response))
}
