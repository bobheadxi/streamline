package streamline_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/jq"
	"go.bobheadxi.dev/streamline/pipeline"
	"go.bobheadxi.dev/streamline/streamexec"
)

func ExampleStream_streamexec() {
	cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")

	stream, _ := streamexec.Start(cmd, streamexec.Combined)
	_ = stream.Stream(func(line string) {
		fmt.Println("received output:", line)
	})
	// Output:
	// received output: hello world
	// received output: this is a line
	// received output: and another line!
}

func ExampleStream_jq() {
	data := strings.NewReader(`Loading...
Still loading...
{
	"message": "this is the real data!"
}`)

	stream := streamline.New(data).
		// Pipeline to discard loading indicators
		WithPipeline(pipeline.Filter(func(line []byte) bool {
			return !bytes.Contains(line, []byte("..."))
		}))

	// stream is just an io.Reader
	message, _ := jq.Query(stream, ".message")

	fmt.Println(string(message))
	// Output: "this is the real data!"
}

func ExampleStream_WithPipeline() {
	data := strings.NewReader("3\n4\n4.8\n7\n5\n2")

	lines, _ := streamline.New(data).
		// Pipeline to discard even and non-integer numbers
		WithPipeline(pipeline.Filter(func(line []byte) bool {
			v, err := strconv.Atoi(string(line))
			return err == nil && v%2 == 0
		})).
		Lines()
	fmt.Println(lines)
	// Output: [4 2]
}

func ExampleStream_WithPipeline_jq() {
	data := strings.NewReader(`Loading...
Still loading...
{"message": "hello"}
{"message":"world"}
{"message":"robert"}`)

	lines, _ := streamline.New(data).
		WithPipeline(pipeline.MultiPipeline{
			// Pipeline to discard non-JSON lines
			pipeline.Filter(func(line []byte) bool {
				return bytes.HasPrefix(line, []byte{'{'})
			}),
			// Execute JQ query for each line
			jq.Pipeline(".message"),
		}).
		Lines()
	fmt.Println(lines)
	// Output: ["hello" "world" "robert"]
}
