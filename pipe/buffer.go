package pipe

import "github.com/djherbis/buffer"

// MemoryBufferSize denotes the maximum in-memory size of the unbounded buffers created by
// this package. Overflows are written to disk at increments of size FileBuffersSize.
var MemoryBufferSize int64 = 128 * 1024

// FileBuffersSize denotes the size of files created to store buffer overflows after the
// in-memory capacity, MemoryBufferSize, is reached.
var FileBuffersSize = MemoryBufferSize / int64(4)

// makeUnboundedBuffer creates a buffer that create files of size fileBuffersSize after
// the in-memory capacity fills up to store overflow.
func makeUnboundedBuffer() buffer.Buffer {
	return buffer.NewUnboundedBuffer(MemoryBufferSize, FileBuffersSize)
}
