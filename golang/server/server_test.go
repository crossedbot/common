package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonResponse(t *testing.T) {
	data := Error{Code: ErrNotAllowedCode, Message: "some message"}
	b, err := json.Marshal(data)
	require.Nil(t, err)
	expectedDataString := string(b)
	expectedCode := http.StatusInternalServerError
	rr := httptest.NewRecorder()
	JsonResponse(rr, data, expectedCode)
	require.Equal(t, expectedDataString, rr.Body.String())
	require.Equal(t, expectedCode, rr.Code)
}

func TestCleanPath(t *testing.T) {
	p := "/"
	require.Equal(t, p, cleanPath(p))
	p = "hello/world/"
	require.Equal(t, "hello/world", cleanPath(p))
}
