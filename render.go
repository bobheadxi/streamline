package streamline

import "errors"

var ErrSkipResult = errors.New("skip result")

type LineRenderer[T any] func(line []byte) (T, error)

type Renderer[T any] struct {
	stream     *Stream
	renderLine LineRenderer[T]
}

func Render[T any](stream *Stream, renderer LineRenderer[T]) *Renderer[T] {
	return &Renderer[T]{
		stream:     stream,
		renderLine: renderer,
	}
}

func (r *Renderer[T]) StreamResults(handle func(T) error) error {
	return r.stream.StreamBytes(func(line []byte) error {
		v, err := r.renderLine(line)
		if err != nil {
			if errors.Is(err, ErrSkipResult) {
				return nil // skip without error
			}
			return err
		}
		return handle(v)
	})
}

func (r *Renderer[T]) Results() ([]T, error) {
	results := make([]T, 0, 10)
	return results, r.stream.StreamBytes(func(line []byte) error {
		v, err := r.renderLine(line)
		if err != nil {
			if errors.Is(err, ErrSkipResult) {
				return nil // skip without error
			}
			return err
		}
		results = append(results, v)
		return nil
	})
}
