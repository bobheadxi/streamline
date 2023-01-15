package streamline_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestStream(t *testing.T) {
	newStream := func() *streamline.Stream {
		return streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world"))
	}

	// These test cases should mirror TestStreamWithPipeline
	for _, tc := range []struct {
		generate func(s *streamline.Stream) (any, error)
		wantErr  bool
		want     autogold.Value
	}{
		{
			generate: func(s *streamline.Stream) (any, error) { return s.Lines() },
			want:     autogold.Want("Lines", []string{"foo bar baz", "baz bar", "hello world"}),
		},
		{
			generate: func(s *streamline.Stream) (any, error) { return s.String() },
			want:     autogold.Want("String", "foo bar baz\nbaz bar\nhello world"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				v, err := s.Bytes()
				return string(v), err
			},
			want: autogold.Want("Bytes", "foo bar baz\nbaz bar\nhello world"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				v, err := io.ReadAll(s)
				return string(v), err
			},
			want: autogold.Want("ReadAll", "foo bar baz\nbaz bar\nhello world"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want: autogold.Want("Copy", "foo bar baz\nbaz bar\nhello world"),
		},
	} {
		t.Run(tc.want.Name(), func(t *testing.T) {
			got, err := tc.generate(newStream())
			if tc.wantErr {
				require.Error(t, err)
				tc.want.Equal(t, err.Error())
			} else {
				assert.NoError(t, err)
				tc.want.Equal(t, got)
			}
		})
	}
}

func TestStreamWithPipeline(t *testing.T) {
	// These test cases should mirror TestStream
	newStream := func() *streamline.Stream {
		return streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))
	}

	for _, tc := range []struct {
		generate func(s *streamline.Stream) (any, error)
		wantErr  bool
		want     autogold.Value
	}{
		{
			generate: func(s *streamline.Stream) (any, error) { return s.Lines() },
			want:     autogold.Want("WithPipeline Lines", []string{"foo-bar-baz", "baz-bar", "hello-world"}),
		},
		{
			generate: func(s *streamline.Stream) (any, error) { return s.String() },
			want:     autogold.Want("WithPipeline String", "foo-bar-baz\nbaz-bar\nhello-world"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				v, err := s.Bytes()
				return string(v), err
			},
			want: autogold.Want("WithPipeline Bytes", "foo-bar-baz\nbaz-bar\nhello-world"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				v, err := io.ReadAll(s)
				return string(v), err
			},
			want: autogold.Want("WithPipeline ReadAll", "foo-bar-baz\nbaz-bar\nhello-world\n"),
		},
		{
			generate: func(s *streamline.Stream) (any, error) {
				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want: autogold.Want("WithPipeline Copy", "foo-bar-baz\nbaz-bar\nhello-world\n"),
		},
	} {
		t.Run(tc.want.Name(), func(t *testing.T) {
			got, err := tc.generate(newStream())
			if tc.wantErr {
				require.Error(t, err)
				tc.want.Equal(t, err.Error())
			} else {
				assert.NoError(t, err)
				tc.want.Equal(t, got)
			}
		})
	}
}

func TestStreamReader(t *testing.T) {
	t.Run("no Pipeline", func(t *testing.T) {
		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))

		p := make([]byte, 5)
		n, err := stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("Read: chunk 1", "foo-b").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("Read: chunk 2", "ar-ba").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("Read: chunk 3", "z\n").Equal(t, string(p[:n]))

		all, err := io.ReadAll(stream)
		assert.NoError(t, err)
		autogold.Want("Read: all remaining", "baz-bar\nhello-world\n").Equal(t, string(all))
	})

	t.Run("WithPipeline", func(t *testing.T) {
		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))

		p := make([]byte, 5)
		n, err := stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("WithPipeline.Read: chunk 1", "foo-b").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("WithPipeline.Read: chunk 2", "ar-ba").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("WithPipeline.Read: chunk 3", "z\n").Equal(t, string(p[:n]))

		all, err := io.ReadAll(stream)
		assert.NoError(t, err)
		autogold.Want("WithPipeline.Read: all remaining", "baz-bar\nhello-world\n").Equal(t, string(all))
	})
}
