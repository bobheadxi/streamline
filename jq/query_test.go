package jq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bobheadxi.dev/streamline"
)

func TestQuery(t *testing.T) {
	t.Run("invalid query", func(t *testing.T) {
		s := streamline.New(strings.NewReader(`{
			"foo":"bar"
		}`))
		res, err := Query(s, ".foo{")
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("ok", func(t *testing.T) {
		s := streamline.New(strings.NewReader(`{
			"foo":"bar"
		}`))
		res, err := Query(s, ".foo")
		assert.NoError(t, err)
		assert.Equal(t, `"bar"`, string(res))
	})
}
