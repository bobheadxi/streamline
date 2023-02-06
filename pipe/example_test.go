package pipe_test

import (
	"errors"
	"fmt"

	"go.bobheadxi.dev/streamline/pipe"
)

func ExampleNewStream() {
	writer, stream := pipe.NewStream()
	go func() {
		writer.Write([]byte("some goroutine emitting data\n"))
		for _, v := range []string{"1", "2", "3", "4"} {
			writer.Write([]byte(v + "\n"))
		}
		writer.CloseWithError(errors.New("oh no!"))
	}()

	err := stream.Stream(func(line string) { fmt.Println(line) })
	fmt.Println("propagated error:", err.Error())
	// Output:
	// some goroutine emitting data
	// 1
	// 2
	// 3
	// 4
	// propagated error: oh no!
}
