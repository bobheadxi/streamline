package streamrender

import (
	"go.bobheadxi.dev/streamline"
)

type LineRenderer[T any] func([]byte) (T, error)

type Renderer[T any] struct {
	stream     *streamline.Stream
	renderLine LineRenderer[T]
}

func New[T any](stream *streamline.Stream, renderer LineRenderer[T]) *Renderer[T] {
	return &Renderer[T]{
		stream:     stream,
		renderLine: renderer,
	}
}

func (r *Renderer[T]) Stream(handle func(T) error) error {
	return r.stream.StreamBytes(func(line []byte) error {
		v, err := r.renderLine(line)
		if err != nil {
			return err
		}
		return handle(v)
	})
}

func (r *Renderer[T]) Lines() ([]T, error) {
	results := make([]T, 0, 10)
	return results, r.stream.StreamBytes(func(line []byte) error {
		v, err := r.renderLine(line)
		if err != nil {
			return err
		}
		results = append(results, v)
		return nil
	})
}
