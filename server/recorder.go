package server

import (
	"net/http"
	"time"
)

// Recorder is an implementation of ResponseWriter that remembers the status code that was set so that it can be accessed.
type Recorder struct {
	http.ResponseWriter
	StartTime  time.Time
	StatusCode int
}

// NewRecorder creates and initializes a new recorder.
func NewRecorder(w http.ResponseWriter) *Recorder {
	return &Recorder{
		ResponseWriter: w,
		StartTime:      time.Now(),
	}
}

// WriteHeader invokes the ResponseWriter implementation and stores the status code.
func (r *Recorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.StatusCode = statusCode
}

// Elapsed determines how long the method took to execute.
func (r *Recorder) Elapsed() time.Duration {
	return time.Now().Sub(r.StartTime)
}
