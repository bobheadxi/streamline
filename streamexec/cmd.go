package streamexec

import (
	"os/exec"

	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipe"
	"go.bobheadxi.dev/streamline/pipeline"
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
)

// Cmd is a command with a streamline.Stream attached to it.
type Cmd struct {
	cmd          *exec.Cmd
	stream       *streamline.Stream
	streamWriter pipe.WriterErrorCloser
}

// Run attaches a streamline.Stream to the command, returning a wrapped Cmd that can be
// configured with pipeline.Pipeline and run with (*Cmd).Run(...).
//
// Output piping is handled by buffers created by streamline/pipe.NewStream(...).
func Attach(cmd *exec.Cmd, mode StreamMode) *Cmd {
	streamWriter, stream := pipe.NewStream()

	if mode&Stdout != 0 {
		cmd.Stdout = streamWriter
	}
	if mode&Stderr != 0 {
		cmd.Stderr = streamWriter
	}

	return &Cmd{
		cmd:          cmd,
		stream:       stream,
		streamWriter: streamWriter,
	}
}

// WithPipeline configures Cmd's streamline.Stream with the given pipeline.
func (c *Cmd) WithPipeline(p pipeline.Pipeline) *Cmd {
	c.stream.WithPipeline(p)
	return c
}

// Start starts a command and returns an error if the command fails to start. It also
// starts a goroutine that waits for command completion and stops the pipe appropriately.
//
// It always returns a valid Stream that can be used to collect output from the underlying
// command.
func (c *Cmd) Start() (*streamline.Stream, error) {
	if err := c.cmd.Start(); err != nil {
		return c.stream, err
	}

	go func() {
		err := c.cmd.Wait()
		c.streamWriter.CloseWithError(err)
	}()

	return c.stream, nil
}
