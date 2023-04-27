package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tastycrayon/blog-backend/db"
)

type contextKey struct {
	name string
}

var cookieAccessKeyCtx = &contextKey{name: "CookieManager"}

var SecretKey = func() []byte {
	key := []byte(os.Getenv("SECRET_KEY"))
	if len(key) == 0 {
		key = []byte("SECRET_KEY")
	}
	return key
}()

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(id int64, t time.Duration) (*string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["sub"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(t).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return nil, 0, err
	}
	return &tokenString, time.Now().Add(t).Unix(), nil
}

func GenerateAccessToken(id int64) (*string, int64, error) {
	tokenLastsFor := time.Hour * 1
	return GenerateToken(id, tokenLastsFor)
}
func GenerateRefreshToken(id int64) (*string, int64, error) {
	tokenLastsFor := time.Hour * 24 * 7
	return GenerateToken(id, tokenLastsFor)
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (*int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["sub"] != nil {
		userId := int64(claims["sub"].(float64))
		return &userId, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// cookie
type CookieManager struct {
	Writer       http.ResponseWriter
	UserId       int64
	IsLoggedIn   bool
	RefreshToken string
}

const AccessToken = "access_token"
const RefreshToken = "refresh_token"

// method to write cookie
func (c *CookieManager) SetToken(token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshToken,
		Value:    token,
		HttpOnly: true,
		Path:     "/refresh",
		// SameSite: http.SameSiteLaxMode,
		Expires: time.Now().Add(time.Hour * 24 * 7), //->->-> set proper date later
	})
}

func GetCookieAccess(ctx context.Context) *CookieManager {
	raw, _ := ctx.Value(cookieAccessKeyCtx).(*CookieManager)
	return raw
}

// Middleware decodes the share session cookie and packs the session into context
func AuthMiddleware(db db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// cookie manager
			cookieManager := &CookieManager{
				Writer: w, IsLoggedIn: false, UserId: 0,
			}
			// put it in context
			ctx := context.WithValue(r.Context(), cookieAccessKeyCtx, cookieManager)
			// and call the next with our new context
			r = r.WithContext(ctx)
			// cookie manager end

			header := r.Header.Get("Authorization")
			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
			// validate token
			var tokenStr string
			if strs := strings.Split(header, " "); len(strs) != 2 {
				next.ServeHTTP(w, r)
				return
			} else {
				tokenStr = strs[1]
			}
			//validate jwt token
			userId, err := ParseToken(tokenStr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			refreshToken, err := r.Cookie("refresh_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			// create user and check if user exists in db
			// no need
			cookieManager.UserId = *userId
			//set logged in to true
			cookieManager.IsLoggedIn = true

			cookieManager.RefreshToken = refreshToken.Value
			next.ServeHTTP(w, r)
		})
	}
}
