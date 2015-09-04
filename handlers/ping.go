// Copyright 2015 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import "net/http"

var pong = []byte("pong")

// HTTP handler for /ping.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write(pong)
}
