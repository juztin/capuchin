// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"
	"time"

	"code.minty.io/marbles/encoders/jsonxml"
)

type timestamp struct {
	UTC      time.Time `json:"utc"`
	Local    time.Time `json:"local"`
	EPOCH    int64     `json:"epoch"`
	Offset   float64   `json:"offset"`
	Timezone string    `json:"zone"`
}

// HTTP handler for /time.
func Time(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	zone, offset := now.Zone()
	t := timestamp{
		UTC:      now.UTC(),
		Local:    now.Local(),
		EPOCH:    now.Unix(),
		Offset:   float64(offset) / 60 / 60,
		Timezone: zone,
	}
	jsonxml.Write(w, r, t)
}
