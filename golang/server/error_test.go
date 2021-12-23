package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorString(t *testing.T) {
	code := ErrRequiredParamCode
	msg := "some message"
	expected := fmt.Sprintf("%d: %s", code, msg)
	actual := Error{code, msg}.Error()
	require.Equal(t, expected, actual)
}
