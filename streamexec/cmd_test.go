package streamexec_test

import (
	"os/exec"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/require"

	"go.bobheadxi.dev/streamline/streamexec"
)

func TestRun(t *testing.T) {
	cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")
	stream, err := streamexec.Attach(cmd, streamexec.Combined).Start()
	require.NoError(t, err)

	out, err := stream.String()
	require.NoError(t, err)

	autogold.Want("run output", `hello world
this is a line
and another line!`).Equal(t, out)
}
