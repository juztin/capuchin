package handlers

import (
	"log"
	"net/http"

	"code.minty.io/config"
	"code.minty.io/hancock"
)

var keyFn hancock.KeyFunc

func Init() {
	expires, ok := config.Int("keysExpire")
	if !ok {
		expires = -1
	}
	keyFn = keyFunc(expires)
}

func apiKeys() map[string]string {
	keys := make(map[string]string)

	var m map[string]interface{}
	if o, ok := config.Val("keys"); !ok {
		log.Println("no API keys found in configuration")
	} else if m, ok = o.(map[string]interface{}); !ok {
		log.Println("invalid config section for API keys")
	}

	for k, v := range m {
		if s, ok := v.(string); !ok {
			log.Panicf("invalid API key: %s -> %s", k, v)
		} else {
			keys[k] = s
		}
	}
	return keys
}

func keyFunc(expires int) hancock.KeyFunc {
	keys := apiKeys()
	return func(key string) (string, int) {
		pKey, ok := keys[key]
		if !ok && expires == -2 {
			pKey = "ignore"
		}
		return pKey, expires
	}
}

// Wraps the handler with hancock signing.
func Signed(h http.Handler) http.Handler {
	return hancock.SignedHandler(h, keyFn)
}

// Wraps the func with hancock signing.
func SignedFunc(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return Signed(http.HandlerFunc(fn))
}
