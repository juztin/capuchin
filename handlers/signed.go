// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"log"
	"net/http"

	"code.minty.io/config"
	"code.minty.io/hancock"
)

var KeyFn hancock.KeyFunc

func Init() {
	expires, ok := config.Int("keysExpire")
	if !ok {
		expires = -1
	}
	KeyFn = KeyFunc(expires)
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

func KeyFunc(expires int) hancock.KeyFunc {
	keys := apiKeys()
	return func(key string) (string, int) {
		pKey, ok := keys[key]
		if !ok && expires == -2 {
			pKey = "ignore"
		}
		return pKey, expires
	}
}

// Wraps the handler with hancock signing, using the given LogFunc to pass validation errors to.
func SignedLog(h http.Handler, fn hancock.LogFunc) http.Handler {
	return hancock.SignedHandler(h, KeyFn, fn)
}

// Wraps the func with hancock signing, using the given LogFunc to pass validation errors to.
func SignedLogFunc(fn func(http.ResponseWriter, *http.Request), logFn hancock.LogFunc) http.Handler {
	return SignedLog(http.HandlerFunc(fn), logFn)
}

// Wraps the handler with hancock signing.
func Signed(h http.Handler) http.Handler {
	return SignedLog(h, log.Println)
}

// Wraps the func with hancock signing.
func SignedFunc(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return Signed(http.HandlerFunc(fn))
}
