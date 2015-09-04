// Copyright 2015 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"code.minty.io/marbles/encoders/jsonxml"
)

var jobs []JobStatus

// JobStatus holds information to generate status information.
type JobStatus struct {
	Name    string
	Timeout time.Duration
	Fn      func() error
}

type status struct {
	Name     string  `json:"endpoint"`
	Duration float64 `json:"duration"`
	Error    string  `json:"error,omitempty"`
}

type statuses struct {
	Endpoints []status `json:"endpoints"`
}

// Adds a jobs to the collection for reteiving statuses.
func AddStatusJob(name string, timeout time.Duration, fn func() error) {
	j := JobStatus{name, timeout, fn}
	jobs = append(jobs, j)
}

// Status HTTP endpoint
func Status(w http.ResponseWriter, r *http.Request) {
	jsonxml.Write(w, r, getStatuses(jobs...))
}

func newStatus(name string, start time.Time, err error) (s status) {
	s.Name = name
	if err != nil {
		s.Duration = -1.0
		s.Error = err.Error()
		return
	}
	s.Duration = float64(time.Since(start)) / float64(time.Millisecond)
	return
}

func getStatuses(jobs ...JobStatus) statuses {
	statuses := statuses{}
	finished := make(chan status)

	for _, j := range jobs {
		go func(job JobStatus) {
			finished <- statusFor(job.Name, job.Timeout, job.Fn)
		}(j)
	}

	var count = len(jobs)
	for {
		if count <= 0 {
			break
		}
		select {
		case s := <-finished:
			statuses.Endpoints = append(statuses.Endpoints, s)
			count--
		}
	}

	return statuses
}

func statusFor(name string, timeout time.Duration, fn func() error) status {
	start := time.Now()
	c := make(chan error)
	go func() { c <- fn() }()
	select {
	case err := <-c:
		return newStatus(name, start, err)
	case <-time.After(timeout):
		return newStatus(name, start, fmt.Errorf("timeout exceeded"))
	}
}
