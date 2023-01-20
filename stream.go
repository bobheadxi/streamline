package streamline

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"go.bobheadxi.dev/streamline/pipeline"
)

// Stream enables live, line-by-line manipulation and handling of data through
// (*Stream).WithPipeline(...) and Stream's various aggregation methods. Stream also
// supports standard library features like io.Copy and io.ReadAll.
//
// Stream's aggregation methods ((*Stream).Stream(...), (*Stream).Lines(...), etc) may
// only be used once. Incremental consumers like (*Stream).Read(...) may need to be called
// multiple times to consume all data but should not be used in conjunction with other
// methods.
type Stream struct {
	// reader carries the input data and the current read state.
	reader *bufio.Reader

	// pipeline, if active, must be used to pre-process lines.
	pipeline pipeline.MultiPipeline

	// readBuffer is set by incremental consumers like Read to store.
	readBuffer *bytes.Buffer

	// lineSeparator is used as the read delimiter.
	lineSeparator byte
}

// New creates a Stream that consumes, processes, and emits data from the input.
func New(input io.Reader) *Stream {
	return &Stream{
		reader:        bufio.NewReader(input),
		lineSeparator: '\n',
	}
}

// WithPipeline configures this Stream to process the input data with the given Pipeline
// in all output methods ((*Stream).Stream(...), (*Stream).Lines(...), io.Copy, etc.).
//
// If one or more Pipelines are already configured on this Stream, the given Pipeline
// is applied sequentially after the preconfigured pipelines.
func (s *Stream) WithPipeline(p pipeline.Pipeline) *Stream {
	s.pipeline = append(s.pipeline, p)
	return s
}

// WithLineSeparator configures a custom line separator for this stream. The default is '\n'.
func (s *Stream) WithLineSeparator(seperator byte) *Stream {
	s.lineSeparator = seperator
	return s
}

type LineHandler[T string | []byte] func(line T) error

// Stream passes lines read from the input to the handler as it processes them.
//
// This method will block until the input returns an error. Unless the error is io.EOF,
// it will also propagate the error.
func (s *Stream) Stream(dst LineHandler[string]) error {
	return s.StreamBytes(func(line []byte) error { return dst(string(line)) })
}

// StreamBytes passes lines read from the input to the handler as it processes them.
//
// This method will block until the input returns an error. Unless the error is io.EOF,
// it will also propagate the error.
func (s *Stream) StreamBytes(dst LineHandler[[]byte]) error {
	for {
		_, err := s.readLine(dst)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

// Lines collects all processed output as a slice of strings.
//
// This method will block until the input returns an error. Unless the error is io.EOF,
// it will also propagate the error.
func (s *Stream) Lines() ([]string, error) {
	lines := make([]string, 0, 10)
	return lines, s.Stream(func(line string) error {
		lines = append(lines, line)
		return nil
	})
}

// String collects all processed output as a string.
//
// This method will block until the input returns an error. Unless the error is io.EOF,
// it will also propagate the error.
func (s *Stream) String() (string, error) {
	var sb strings.Builder
	_, err := s.WriteTo(&sb) // use Stream implementation of io
	data := strings.TrimSuffix(sb.String(), "\n")
	return data, err
}

// Bytes collects all processed output as a bytes slice.
//
// This method will block until the input returns an error. Unless the error is io.EOF,
// it will also propagate the error.
func (s *Stream) Bytes() ([]byte, error) {
	var b bytes.Buffer
	_, err := s.WriteTo(&b)
	data := bytes.TrimSuffix(b.Bytes(), []byte{s.lineSeparator})
	return data, err
}

var _ io.WriterTo = (*Stream)(nil)

// WriteTo writes processed data to dst. It allows Stream to effectively implement io.Copy
// handling.
func (s *Stream) WriteTo(dst io.Writer) (int64, error) {
	if s.pipeline.Inactive() {
		// Happy path, do a straight read from the underlying source
		return s.reader.WriteTo(dst)
	}

	var totalWritten int64
	return totalWritten, s.StreamBytes(func(line []byte) error {
		n, err := dst.Write(append(line, '\n'))
		totalWritten += int64(n)
		return err
	})
}

var _ io.Reader = (*Stream)(nil)

// Read populates p with processed data. It allows Stream to effectively be compatible
// with anything that accepts an io.Reader.
func (s *Stream) Read(p []byte) (int, error) {
	if s.pipeline.Inactive() {
		// Happy path, do a straight read
		return s.reader.Read(p)
	}

	if s.readBuffer == nil {
		s.readBuffer = &bytes.Buffer{}
	}
	if s.readBuffer.Len() > 0 {
		read, _ := s.readBuffer.Read(p)
		return read, nil
	}

	// TODO: Going line-by-line may require a lot of calls to Read() in scenarios where
	// each line is very short. This implementation currently only reads multiple lines
	// if we a read is a skip - we may want to implement the ability to read again if p
	// is not filled.
	var (
		read    int
		skipped bool
		err     error
	)
	for {
		skipped, err = s.readLine(func(line []byte) error {
			line = append(line, s.lineSeparator)
			if len(line) <= cap(p) {
				// If we have less data than is requested, just copy into p
				read = copy(p, line)
			} else {
				// If we are using readLine, then the buffer has been consumed,
				// we reset and populate it with our data and give p a partial
				// read.
				s.readBuffer.Reset()
				_, _ = s.readBuffer.Write(line)
				read, _ = s.readBuffer.Read(p)
			}
			return nil
		})
		if !skipped || err != nil {
			break
		}
	}
	return read, err
}

// readLine consumes a single line in the stream. The error returned, in order of
// precedence, is one of:
//
// - processing error
// - handler error
// - read error
//
// The read error in particular may be io.EOF, which the caller should handle on a
// case-by-case basis.
func (s *Stream) readLine(handle LineHandler[[]byte]) (skipped bool, err error) {
	line, readErr := s.reader.ReadBytes(s.lineSeparator)

	if len(line) == 0 {
		return true, readErr
	}

	// If the line ends with a newline, trim it before handing it off to the
	// LineHandler.
	if line[len(line)-1] == s.lineSeparator {
		line = line[:len(line)-1]
	}

	var processErr error
	if line, processErr = s.pipeline.ProcessLine(line); processErr != nil {
		return false, processErr
	}

	if line == nil {
		return true, nil
	}

	if dstErr := handle(line); dstErr != nil {
		return false, dstErr
	}

	return false, readErr
}
