// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Logger func(interface{}) error

// HTTP panic recovery handler.
type Recover struct {
	handler http.Handler
}

var Log Logger = func(o interface{}) error {
	log.Println(o)
	return nil
}

func (h *Recover) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)

			//errMsg := fmt.Sprintf("%s", err)
			//buf := bytes.NewBuffer([]byte(errMsg))
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("%s", err))
			for i := 1; ; i++ {
				if _, f, l, ok := runtime.Caller(i); !ok {
					break
				} else {
					buf.WriteString(fmt.Sprintf("Line: %d\nfile: %s\n-\n", l, f))
				}
			}
			if err := Log(buf.String()); err != nil {
				log.Println(err)
			}
		}
	}()
	h.handler.ServeHTTP(w, r)
}

// Wraps the func with a panic recovery.
func Recovery(h http.Handler) http.Handler {
	return &Recover{h}
}

// Wraps the func with a panic recovery.
func RecoveryFunc(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return Recovery(http.HandlerFunc(fn))
}
