package streamline_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
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
		name     string
		generate func(s *streamline.Stream) (any, error)
		wantErr  bool
		want     autogold.Value
	}{
		{
			name:     "Lines",
			generate: func(s *streamline.Stream) (any, error) { return s.Lines() },
			want:     autogold.Expect([]string{"foo bar baz", "baz bar", "hello world"}),
		},
		{
			name:     "String",
			generate: func(s *streamline.Stream) (any, error) { return s.String() },
			want:     autogold.Expect("foo bar baz\nbaz bar\nhello world"),
		},
		{
			name: "Bytes",
			generate: func(s *streamline.Stream) (any, error) {
				v, err := s.Bytes()
				return string(v), err
			},
			want: autogold.Expect("foo bar baz\nbaz bar\nhello world"),
		},
		{
			name: "io.ReadAll",
			generate: func(s *streamline.Stream) (any, error) {
				v, err := io.ReadAll(s)
				return string(v), err
			},
			want: autogold.Expect("foo bar baz\nbaz bar\nhello world"),
		},
		{
			name: "io.Copy",
			generate: func(s *streamline.Stream) (any, error) {
				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want: autogold.Expect("foo bar baz\nbaz bar\nhello world"),
		},
		{
			name: "StreamBytes",
			generate: func(s *streamline.Stream) (any, error) {
				var count int
				return nil, s.StreamBytes(func(line []byte) error {
					count += 1
					if count > 2 {
						return errors.New("oh no!")
					}
					return nil
				})
			},
			want:    autogold.Expect("oh no!"),
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

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
		name     string
		generate func(s *streamline.Stream) (any, error)
		wantErr  bool
		want     autogold.Value
	}{
		{
			name:     "Lines",
			generate: func(s *streamline.Stream) (any, error) { return s.Lines() },
			want:     autogold.Expect([]string{"foo-bar-baz", "baz-bar", "hello-world"}),
		},
		{
			name:     "String",
			generate: func(s *streamline.Stream) (any, error) { return s.String() },
			want:     autogold.Expect("foo-bar-baz\nbaz-bar\nhello-world"),
		},
		{
			name: "Bytes",
			generate: func(s *streamline.Stream) (any, error) {
				v, err := s.Bytes()
				return string(v), err
			},
			want: autogold.Expect("foo-bar-baz\nbaz-bar\nhello-world"),
		},
		{
			name: "io.ReaadAll",
			generate: func(s *streamline.Stream) (any, error) {
				v, err := io.ReadAll(s)
				return string(v), err
			},
			want: autogold.Expect("foo-bar-baz\nbaz-bar\nhello-world\n"),
		},
		{
			name: "io.Copy",
			generate: func(s *streamline.Stream) (any, error) {
				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want: autogold.Expect("foo-bar-baz\nbaz-bar\nhello-world\n"),
		},
		{
			name: "multiple pipelines",
			generate: func(s *streamline.Stream) (any, error) {
				s = s.WithPipeline(pipeline.Map(func(line []byte) []byte {
					return bytes.ReplaceAll(line, []byte("bar"), []byte("robert"))
				}))

				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want: autogold.Expect("foo-robert-baz\nbaz-robert\nhello-world\n"),
		},
		{
			name: "pipeline error",
			generate: func(s *streamline.Stream) (any, error) {
				var count int
				s = s.WithPipeline(pipeline.MapErr(func(line []byte) ([]byte, error) {
					count += 1
					if count > 2 {
						return nil, errors.New("oh no!")
					}
					return line, nil
				}))

				var sb strings.Builder
				_, err := io.Copy(&sb, s)
				return sb.String(), err
			},
			want:    autogold.Expect("oh no!"),
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

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
		t.Parallel()

		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world"))

		t.Run("chunk 1", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("foo b").Equal(t, string(p[:n]))
		})

		t.Run("chunk 2", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("ar ba").Equal(t, string(p[:n]))
		})

		t.Run("chunk 3", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("z\nbaz").Equal(t, string(p[:n]))
		})

		t.Run("remaining", func(t *testing.T) {
			all, err := io.ReadAll(stream)
			assert.NoError(t, err)
			autogold.Expect(" bar\nhello world").Equal(t, string(all))
		})

		t.Run("read on empty input", func(t *testing.T) {
			all, err := io.ReadAll(stream)
			assert.Zero(t, len(all))
			assert.NoError(t, err)
		})
	})

	t.Run("WithPipeline", func(t *testing.T) {
		t.Parallel()

		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))

		t.Run("chunk 1", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("foo-b").Equal(t, string(p[:n]))
		})

		t.Run("chunk 2", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("ar-ba").Equal(t, string(p[:n]))
		})

		t.Run("chunk 3", func(t *testing.T) {
			p := make([]byte, 5)
			n, err := stream.Read(p)
			assert.NoError(t, err)
			assert.NotZero(t, n)
			autogold.Expect("z\n").Equal(t, string(p[:n]))
		})

		t.Run("remaining", func(t *testing.T) {
			all, err := io.ReadAll(stream)
			assert.NoError(t, err)
			autogold.Expect("baz-bar\nhello-world\n").Equal(t, string(all))
		})

		t.Run("read on empty input", func(t *testing.T) {
			all, err := io.ReadAll(stream)
			assert.Zero(t, len(all))
			assert.NoError(t, err)
		})
	})

	t.Run("with all lines skipped", func(t *testing.T) {
		t.Parallel()

		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return nil
			}))

		// Correctly read nothing
		all, err := io.ReadAll(stream)
		assert.Zero(t, len(all))
		assert.NoError(t, err)
	})

	t.Run("with some lines skipped", func(t *testing.T) {
		t.Parallel()

		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Sample(2)) // only line 2 will stay

		// Correctly read only 2nd line
		all, err := io.ReadAll(stream)
		assert.NoError(t, err)
		autogold.Expect("baz bar\n").Equal(t, string(all))
	})
}

func TestStreamWithLineSeparator(t *testing.T) {
	data := "Compressing objects:   0% (1/4334)\rCompressing objects:   1% (44/4334)\rCompressing objects:   2% (87/4334)\rCompressing objects:   3% (131/4334)\rCompressing objects:   4% (174/4334)\rCompressing objects:   5% (217/4334)\rCompressing objects:   6% (261/4334)\rCompressing objects:   7% (304/4334)\rCompressing objects:   8% (347/4334)\rCompressing objects:   9% (391/4334)\rCompressing objects:  10% (434/4334)\rCompressing objects:  11% (477/4334)\rCompressing objects:  12% (521/4334)\rCompressing objects:  13% (564/4334)\rCompressing objects:  14% (607/4334)\rCompressing objects:  15% (651/4334)\rCompressing objects:  16% (694/4334)\rCompressing objects:  17% (737/4334)"
	stream := streamline.New(strings.NewReader(data)).WithLineSeparator('\r')

	lines, err := stream.Lines()
	require.NoError(t, err)
	autogold.Expect([]string{
		"Compressing objects:   0% (1/4334)", "Compressing objects:   1% (44/4334)",
		"Compressing objects:   2% (87/4334)",
		"Compressing objects:   3% (131/4334)",
		"Compressing objects:   4% (174/4334)",
		"Compressing objects:   5% (217/4334)",
		"Compressing objects:   6% (261/4334)",
		"Compressing objects:   7% (304/4334)",
		"Compressing objects:   8% (347/4334)",
		"Compressing objects:   9% (391/4334)",
		"Compressing objects:  10% (434/4334)",
		"Compressing objects:  11% (477/4334)",
		"Compressing objects:  12% (521/4334)",
		"Compressing objects:  13% (564/4334)",
		"Compressing objects:  14% (607/4334)",
		"Compressing objects:  15% (651/4334)",
		"Compressing objects:  16% (694/4334)",
		"Compressing objects:  17% (737/4334)",
	}).Equal(t, lines)
}
