package streamexec

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipe"
)

// StreamMode indicates what output(s) to attach.
type StreamMode int

const (
	// Combined streams both Stdout and Stderr.
	Combined StreamMode = Stdout | Stderr

	// Stdout only streams cmd.Stdout.
	Stdout StreamMode = 1 << iota
	// Stderr only streams cmd.Stderr.
	Stderr

	// ErrWithStderr collects Stderr output and includes it in the returned error from
	// Cmd.Start(). Best used with the Stdout StreamMode.
	ErrWithStderr
)

// Start attaches a streamline.Stream to the command and starts it. It returns an error
// if the command fails to start. If the command succesfully starts, it also starts a
// goroutine that waits for command completion and stops the pipe appropriately.
//
// Instead of using cmd.Wait() for command completion, callers should read the returned
// Stream until completion to indicate if the command has exited.
//
// Before consuming the Stream, the caller can configure the Stream as a normal stream
// using e.g. WithPipeline.
//
// Output piping is handled by buffers created by streamline/pipe.NewStream(...).
func Start(cmd *exec.Cmd, mode StreamMode) (*streamline.Stream, error) {
	streamWriter, stream := pipe.NewStream()

	if mode&Stdout != 0 {
		cmd.Stdout = streamWriter
	}
	if mode&Stderr != 0 {
		cmd.Stderr = streamWriter
	}

	var stderr *strings.Builder
	if mode&ErrWithStderr != 0 {
		stderr = &strings.Builder{}
		if cmd.Stderr == nil {
			cmd.Stderr = stderr
		} else {
			cmd.Stderr = io.MultiWriter(cmd.Stderr, stderr)
		}
	}

	if err := cmd.Start(); err != nil {
		return stream, err
	}

	go func() {
		err := cmd.Wait()
		if err != nil && stderr != nil && stderr.Len() > 0 {
			err = fmt.Errorf("%w: %s", err, strings.TrimSuffix(stderr.String(), "\n"))
		}
		streamWriter.CloseWithError(err)
	}()

	return stream, nil
}
