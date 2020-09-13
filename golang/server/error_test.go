package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorString(t *testing.T) {
	e := Error{Code: 123, Title: "some title", Text: "some text"}
	expected := fmt.Sprintf("[%d] %s: %s", e.Code, e.Title, e.Text)
	require.Equal(t, expected, e.String())
}
