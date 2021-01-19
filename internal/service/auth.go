package service

import (
	"encoding/base64"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"unsafe"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

// Accounts defines a key/value for user/pass list of authorized logins.
type Accounts map[string]string

type authPair struct {
	value string
	user  string
}

type authPairs []authPair

func basicAuth(accounts gin.Accounts) gin.HandlerFunc {
	realm := "Authorization Required"
	pairs := processAccounts(accounts)
	return func(c *gin.Context) {
		user, found := pairs.searchCredential(c.Request.Header.Get("Authorization"))
		if !found {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(AuthUserKey, user)
	}
}

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if utils.CheckPasswordHash(pair.value, authValue) {
			return pair.user, true
		}
	}
	return "", false
}

func processAccounts(accounts gin.Accounts) authPairs {
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			value: value,
			user:  user,
		})
	}
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(stringToBytes(base))
}

func stringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}