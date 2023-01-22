package pipe

import (
	"io"

	"github.com/djherbis/nio/v3"

	"go.bobheadxi.dev/streamline"
)

type WriterErrorCloser interface {
	io.Writer
	// CloseWithError will prevent further writes to this Stream and propagate the error
	// to Stream readers.
	CloseWithError(error) error
}

// NewStream creates a Stream that consumes and emits data written to the returned writer,
// piped through a preconfigured, unbounded buffer (see MemoryBufferSize and
// FileBuffersSize).
//
// The writer should be closed by the caller when all data has been written or when an
// error occurs at the source.
func NewStream() (writer WriterErrorCloser, stream *streamline.Stream) {
	outputBuffer := makeUnboundedBuffer()

	// We use this buffered pipe from github.com/djherbis/nio that allows async read and
	// write operations to the reader and writer portions of the pipe respectively.
	outputReader, outputWriter := nio.Pipe(outputBuffer)

	return outputWriter, streamline.New(outputReader)
}