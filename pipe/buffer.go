package pipe

import "github.com/djherbis/buffer"

// MemoryBufferSize denotes the maximum in-memory size of the buffers created by this
// package.
//
// When using the default unbounded buffer in NewStream, overflows are written to disk at
// increments of size FileBuffersSize.
var MemoryBufferSize int64 = 128 * 1024

// FileBuffersSize denotes the size of files created to store buffer overflows after the
// in-memory capacity, MemoryBufferSize, is reached in the unbounded buffers created by
// NewStream.
var FileBuffersSize = MemoryBufferSize / int64(4)

// makeUnboundedBuffer creates a buffer that create files of size fileBuffersSize after
// the in-memory capacity fills up to store overflow.
func makeUnboundedBuffer() buffer.Buffer {
	return buffer.NewUnboundedBuffer(MemoryBufferSize, FileBuffersSize)
}

// makeMemoryBuffer creates a buffer that only works up to MemoryBufferSize, and never
// overflows to disk unlike the unbounded buffer.
func makeMemoryBuffer() buffer.Buffer {
	return buffer.New(MemoryBufferSize)
}
