package streamexec_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/require"

	"go.bobheadxi.dev/streamline/pipeline"
	"go.bobheadxi.dev/streamline/streamexec"
)

func TestRun(t *testing.T) {
	cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")
	stream, err := streamexec.Start(cmd) // should default to combined
	require.NoError(t, err)

	out, err := stream.String()
	require.NoError(t, err)

	autogold.Want("run output", `hello world
this is a line
and another line!`).Equal(t, out)

	t.Run("WithPipeline", func(t *testing.T) {
		cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")
		stream, err := streamexec.Start(cmd, streamexec.Combined)
		require.NoError(t, err)

		out, err := stream.
			WithPipeline(pipeline.Filter(func(line []byte) bool {
				return !bytes.Contains(line, []byte("hello"))
			})).
			String()
		require.NoError(t, err)

		autogold.Want("run output with pipeline", "this is a line\nand another line!").Equal(t, out)
	})
}

func TestError(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cmd := exec.Command("bash", "-o", "pipefail", "-o", "errexit", "-c",
			`echo "stdout" ; sleep 0.001 ; >&2 echo "stderr" ; exit 1`)
		stream, err := streamexec.Start(cmd, streamexec.Stdout|streamexec.ErrWithStderr)
		require.NoError(t, err)

		out, err := stream.String()
		require.Error(t, err)
		autogold.Want("got error with stderr", "exit status 1: stderr").Equal(t, err.Error())
		autogold.Want("get only stdout", "stdout").Equal(t, out)
	})
	t.Run("with existing stderr", func(t *testing.T) {
		cmd := exec.Command("bash", "-o", "pipefail", "-o", "errexit", "-c",
			`echo "stdout" ; sleep 0.001 ; >&2 echo "stderr" ; exit 1`)
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		stream, err := streamexec.Start(cmd, streamexec.Stdout, streamexec.ErrWithStderr)
		require.NoError(t, err)

		_, err = stream.String()
		require.Error(t, err)
		autogold.Want("get stderr in assigned buffer", "stderr\n").Equal(t, errBuf.String())
		autogold.Want("got error with stderr as well", "exit status 1: stderr").Equal(t, err.Error())
	})
}
