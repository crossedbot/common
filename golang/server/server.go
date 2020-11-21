package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/crossedbot/common/golang/logger"
)

// Server is interface that represents an HTTP server.
type Server interface {
	Start() error
	Stop() error
	Reload() error
	Add(handler Handler, method, path string, settings ...ResponseSetting) error
}

// server implements the Server interface.
type server struct {
	addr string             // server address
	rto  int                // reader timeout
	rtr  *httprouter.Router // router
	run  int32              // indicates whether the server is running or not atomically
	srv  *http.Server       // server
	wg   sync.WaitGroup     // tracks pending requests
	wto  int                // writer timeout
}

// New returns a server at the given address.
func New(addr string, readTimeoutSeconds, writeTimeoutSeconds int) Server {
	return &server{
		addr: addr,
		rto:  readTimeoutSeconds,
		rtr:  router(),
		wto:  writeTimeoutSeconds,
	}
}

// Start starts the server for accepting requests.
func (s *server) Start() error {
	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      s.rtr,
		ReadTimeout:  time.Duration(s.rto) * time.Second,
		WriteTimeout: time.Duration(s.wto) * time.Second,
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to create listener; %s", err.Error())
	}
	go s.srv.Serve(listener)
	atomic.StoreInt32(&s.run, 1)
	return nil
}

// Stop stops the server from accepting requests.
func (s *server) Stop() error {
	atomic.StoreInt32(&s.run, 0)
	s.wg.Wait()
	return s.srv.Shutdown(context.Background())
}

// Reload restarts the server.
func (s *server) Reload() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

// Add adds a new handler for the given method at the given path; applying
// all response settings to the response writer. Allowable methods include: GET,
// HEAD, POST, PUT, PATCH, DELETE, and OPTIONS. CONNECT and TRACE are not
// supported.
func (s *server) Add(handler Handler, method, path string, settings ...ResponseSetting) error {
	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s.wg.Add(1)
		defer s.wg.Done()
		if atomic.LoadInt32(&s.run) < 1 {
			JsonResponse(
				w,
				ErrorServiceUnavailable,
				http.StatusServiceUnavailable,
			)
		}
		applyResponseSettings(w, settings)
		handler(w, r, parameters(p))
	}
	path = cleanPath(path)
	switch method {
	case http.MethodGet:
		s.rtr.GET(path, h)
	case http.MethodHead:
		s.rtr.HEAD(path, h)
	case http.MethodPost:
		s.rtr.POST(path, h)
	case http.MethodPut:
		s.rtr.PUT(path, h)
	case http.MethodPatch:
		s.rtr.PATCH(path, h)
	case http.MethodDelete:
		s.rtr.DELETE(path, h)
	case http.MethodOptions:
		s.rtr.OPTIONS(path, h)
	default:
		return fmt.Errorf("method %s is not supported", method)
	}
	return nil
}

// JsonResponse encodes and writes a JSON response using the given data object.
func JsonResponse(w http.ResponseWriter, data interface{}, status int) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", b)
}

// router sets up a new httprouter.Router with predefined handlers for panics,
// resource not found, and method not allowed.
func router() *httprouter.Router {
	rtr := httprouter.New()
	rtr.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		logger.Error(err)
		JsonResponse(w, ErrProcessingRequest, http.StatusInternalServerError)
	}
	rtr.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JsonResponse(w, ErrNotFound, http.StatusInternalServerError)
	})
	rtr.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JsonResponse(w, ErrNotAllowed, http.StatusInternalServerError)
	})
	return rtr
}

// cleanPath is a utility function to clean up a url path.
func cleanPath(p string) string {
	if p == "/" {
		return p
	}
	return strings.TrimSuffix(p, "/")
}
