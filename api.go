package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

// api server handler func type
type APIFunc func(w http.ResponseWriter, r *http.Request) error

// API serve options
type APIServerOpts struct {
	ListenAddr string
	localNode  *FileServer
}

// api server
type APIServer struct {
	APIServerOpts
	mux *http.ServeMux
}

// New takes APIServerOpts and creates a new APIServer
// returns new *APIServer
func NewAPIServer(opts APIServerOpts) *APIServer {
	return &APIServer{
		mux:           http.NewServeMux(),
		APIServerOpts: opts,
	}
}

func (s *APIServer) Run() error {

	// TODO: upload files
	s.mux.HandleFunc("POST /upload", s.handler(s.write))
	// TODO: retrieve files

	// start and listen api server
	return http.ListenAndServe(s.ListenAddr, s.mux)
}

func (s *APIServer) handler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create new context from request context with timeout of 30 seconds
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// run api func
		err := fn(w, r.WithContext(ctx))
		if err != nil {
			// handle error
			log.Panic(err)
		}
	}
}

func (s *APIServer) write(w http.ResponseWriter, r *http.Request) error {

	s.localNode.Store("12345", r.Body)

	return nil
}
