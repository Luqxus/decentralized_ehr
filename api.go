package main

import (
	"context"
	"encoding/json"
	"io"
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
	s.mux.HandleFunc("POST /write", s.handler(s.write))

	// TODO: retrieve files
	s.mux.HandleFunc("GET /read", s.handler(s.read))

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

	err := s.localNode.Store("12345", r.Body)
	if err != nil {
		http.Error(w, "failed to write data", http.StatusInternalServerError)
		return nil
	}

	return writeJSON(w, map[string]string{"message": "data written successfully"})
}

func (s *APIServer) read(w http.ResponseWriter, r *http.Request) error {
	key, ok := r.URL.Query()["key"]
	if !ok {
		http.Error(w, "key not found", http.StatusBadRequest)
		return nil
	}

	n, reader, err := s.localNode.Get(key[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	b := make([]byte, n)
	_, err = reader.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	_, err = w.Write(b)
	return err
}

func writeJSON(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
