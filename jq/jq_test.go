package jq

import (
	"strings"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
)

func TestPipeline(t *testing.T) {
	t.Run("invalid query", func(t *testing.T) {
		p := Pipeline("asdf{")
		assert.False(t, p.Inactive())

		l, err := p.ProcessLine([]byte(`{"foo":"bar"}`))
		assert.Empty(t, l)
		require.Error(t, err)
		autogold.Want("invalid query error", `jq.Parse: unexpected token "{"`).Equal(t, err.Error())
	})

	t.Run("query error", func(t *testing.T) {
		p := Pipeline(".baz[4]")
		assert.False(t, p.Inactive())

		l, err := p.ProcessLine([]byte(`{"foo":bar}`))
		assert.Empty(t, l)
		require.Error(t, err)
		autogold.Want("query failed error", `json: invalid character 'b' looking for beginning of value: {"foo":bar}`).Equal(t, err.Error())
	})

	t.Run("ok", func(t *testing.T) {
		p := Pipeline(".foo")
		assert.False(t, p.Inactive())

		l, err := p.ProcessLine([]byte(`{"foo":"bar"}`))
		assert.NoError(t, err)
		assert.Equal(t, `"bar"`, string(l))
	})
}

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
