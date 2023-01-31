package jq_test

import (
	"bytes"
	"fmt"
	"strings"

	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/jq"
	"go.bobheadxi.dev/streamline/pipeline"
)

func ExamplePipeline() {
	data := strings.NewReader(`{"message": "hello"}
{"message":"world"}
{"message":"robert"}`)

	lines, err := streamline.New(data).
		WithPipeline(jq.Pipeline(".message")).
		Lines()
	if err != nil {
		fmt.Println("stream failed:", err.Error())
	}
	fmt.Println(lines)
	// Output: ["hello" "world" "robert"]
}

func ExampleQuery() {
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
	message, err := jq.Query(stream, ".message")
	if err != nil {
		fmt.Println("query failed:", err.Error())
	}

	fmt.Println(string(message))
	// Output: "this is the real data!"
}
