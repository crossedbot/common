package server

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler represents an HTTP handler method.
type Handler func(http.ResponseWriter, *http.Request, Parameters)

func NewHandler(handler http.HandlerFunc) Handler {
	return func(w http.ResponseWriter, r *http.Request, p Parameters) {
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, p.base())
		handler(w, r.WithContext(ctx))
	}
}

// base returns the Handler object as a httprouter.Handle object.
func (h Handler) base() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ps := parameters(p)
		h(w, r, ps)
	}
}

// Parameter represents a single URL parameter's key-value pair.
type Parameter struct {
	Key   string
	Value string
}

// Parameters represents a list of URL parameters.
type Parameters []Parameter

func GetParameters(ctx context.Context) Parameters {
	return parameters(httprouter.ParamsFromContext(ctx))
}

// Get returns a parmeter value for the given key. If a key does not exist an
// empty string is returned.
func (params Parameters) Get(key string) string {
	for _, p := range params {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// base returns the Parameters object as a httprouter.Params object.
func (p Parameters) base() httprouter.Params {
	params := httprouter.Params{}
	for _, v := range p {
		params = append(params, httprouter.Param(v))
	}
	return params
}

// parameters converts a httprouter.Params object into a Parameters object.
func parameters(p httprouter.Params) Parameters {
	params := Parameters{}
	for _, v := range p {
		params = append(params, Parameter(v))
	}
	return params
}

// ResponseSetting can be used to configure a repsonse writer.
type ResponseSetting func(w *http.ResponseWriter)

// SetResponseHeader is a header setting that sets the response header key-value
// pair.
func SetResponseHeader(key, value string) ResponseSetting {
	return func(w *http.ResponseWriter) {
		(*w).Header().Set(key, value)
	}
}

// applyResponseSettings applies all settings to a response writer.
func applyResponseSettings(w *http.ResponseWriter, settings []ResponseSetting) {
	for _, s := range settings {
		s(w)
	}
}
