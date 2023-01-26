package pipe

import (
	"io"

	// We use this buffered pipe from github.com/djherbis/nio that allows async read and
	// write operations to the reader and writer portions of the pipe respectively.
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
// FileBuffersSize). The writer and reader portions of the pipe can be written to and read
// asynchronously.
//
// The returned WriterErrorCloser must be closed by the caller when all data has been
// written or when an error occurs at the source to indicate to the stream that no further
// data will become available.
func NewStream() (writer WriterErrorCloser, stream *streamline.Stream) {
	outputBuffer := makeUnboundedBuffer()

	outputReader, outputWriter := nio.Pipe(outputBuffer)

	return outputWriter, streamline.New(outputReader)
}

// NewStream creates a Stream that consumes and emits data written to the returned writer,
// piped through a preconfigured bounded buffer (see MemoryBufferSize). It may be
// preferred over the default NewStream if you do not want to buffer to overflow onto
// temporary files.  The writer and reader portions of the pipe can be written to and read
// asynchronously.
//
// The returned WriterErrorCloser must be closed by the caller when all data has been
// written or when an error occurs at the source to indicate to the stream that no further
// data will become available.
func NewBoundedStream() (writer WriterErrorCloser, stream *streamline.Stream) {
	outputBuffer := makeMemoryBuffer()

	outputReader, outputWriter := nio.Pipe(outputBuffer)

	return outputWriter, streamline.New(outputReader)
}
