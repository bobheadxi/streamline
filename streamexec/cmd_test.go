package streamexec_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bobheadxi.dev/streamline/pipeline"
	"go.bobheadxi.dev/streamline/streamexec"
)

func TestRun(t *testing.T) {
	t.Run("default to Combined output", func(t *testing.T) {
		t.Parallel()

		cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")
		stream, err := streamexec.Start(cmd) // should default to combined
		require.NoError(t, err)

		out, err := stream.String()
		require.NoError(t, err)

		autogold.Expect(`hello world
this is a line
and another line!`).Equal(t, out)
	})

	t.Run("WithPipeline", func(t *testing.T) {
		t.Parallel()

		cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")
		stream, err := streamexec.Start(cmd, streamexec.Combined)
		require.NoError(t, err)

		out, err := stream.
			WithPipeline(pipeline.Filter(func(line []byte) bool {
				return !bytes.Contains(line, []byte("hello"))
			})).
			String()
		require.NoError(t, err)

		autogold.Expect("this is a line\nand another line!").Equal(t, out)
	})

	t.Run("stderr and exit", func(t *testing.T) {
		t.Parallel()

		cmd := exec.Command("bash", "-o", "pipefail", "-o", "errexit", "-c",
			`echo "stdout" ; sleep 0.001 ; >&2 echo "stderr" ; exit 1`)
		stream, err := streamexec.Start(cmd, streamexec.Stdout|streamexec.ErrWithStderr)
		require.NoError(t, err)

		out, err := stream.String()
		require.Error(t, err)
		autogold.Expect("exit status 1: stderr").Equal(t, err.Error())
		autogold.Expect("stdout").Equal(t, out)
	})

	t.Run("multiple stderr and exit", func(t *testing.T) {
		t.Parallel()

		cmd := exec.Command("bash", "-o", "pipefail", "-o", "errexit", "-c",
			`echo "stdout" ; sleep 0.001 ; >&2 echo "stderr" ; exit 1`)
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		stream, err := streamexec.Start(cmd, streamexec.Stdout, streamexec.ErrWithStderr)
		require.NoError(t, err)

		_, err = stream.String()
		require.Error(t, err)
		autogold.Expect("stderr\n").Equal(t, errBuf.String())
		autogold.Expect("exit status 1: stderr").Equal(t, err.Error())
	})

	t.Run("failed to start", func(t *testing.T) {
		t.Parallel()

		cmd := exec.Command("foobar")
		stream, err := streamexec.Start(cmd)
		assert.Error(t, err)
		assert.NotNil(t, stream)

		out, err := stream.String()
		assert.NoError(t, err)
		assert.Empty(t, out)
	})
}
