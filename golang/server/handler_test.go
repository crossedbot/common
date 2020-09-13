package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedBody := "Hello World"
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
		fmt.Fprintf(w, expectedBody)
	}
	handler := NewHandler(handlerFunc)
	r, err := http.NewRequest(http.MethodGet, "/hello", nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	handler(rr, r, Parameters{})
	require.Equal(t, expectedStatus, rr.Code)
	require.Equal(t, expectedBody, rr.Body.String())
}

func TestHandlerBase(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedBody := "Hello World"
	baseHandler := Handler(
		func(w http.ResponseWriter, r *http.Request, p Parameters) {
			w.WriteHeader(expectedStatus)
			fmt.Fprintf(w, expectedBody)
		},
	).base()
	r, err := http.NewRequest(http.MethodGet, "/hello", nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	baseHandler(rr, r, httprouter.Params{})
	require.Equal(t, expectedStatus, rr.Code)
	require.Equal(t, expectedBody, rr.Body.String())
}

func TestGetParameters(t *testing.T) {
	expected := Parameters{
		Parameter{
			Key:   "id",
			Value: "abc",
		},
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, expected.base())
	actual := GetParameters(ctx)
	require.Equal(t, expected, actual)
}

func TestParametersGet(t *testing.T) {
	key := "id"
	value := "abc"
	p := Parameters{Parameter{Key: key, Value: value}}
	require.Equal(t, value, p.Get(key))
}

func TestParametersBase(t *testing.T) {
	key := "id"
	value := "abc"
	p := Parameters{Parameter{Key: key, Value: value}}
	baseP := p.base()
	require.Equal(t, value, baseP.ByName(key))
}

func TestParameters(t *testing.T) {
	key := "id"
	value := "abc"
	p := httprouter.Params{httprouter.Param{Key: key, Value: value}}
	require.Equal(t, value, parameters(p).Get(key))
}

func TestResponseHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	key := "auth"
	value := "bearer"
	setting := SetResponseHeader(key, value)
	setting(rr)
	require.Equal(t, value, rr.Header().Get(key))
}

func TestApplyResponseSettings(t *testing.T) {
	rr := httptest.NewRecorder()
	key := "auth"
	value := "bearer"
	setting := SetResponseHeader(key, value)
	applyResponseSettings(rr, []ResponseSetting{setting})
	require.Equal(t, value, rr.Header().Get(key))
}
