package streamexec_test

import (
	"fmt"
	"os/exec"

	"go.bobheadxi.dev/streamline/streamexec"
)

func ExampleStart() {
	cmd := exec.Command("echo", "hello world\nthis is a line\nand another line!")

	stream, err := streamexec.Start(cmd, streamexec.Combined)
	if err != nil {
		fmt.Println("failed to start:", err.Error())
	}
	err = stream.Stream(func(line string) {
		fmt.Println("received output:", line)
	})
	if err != nil {
		fmt.Println("command exited with error status:", err.Error())
	}
	// Output:
	// received output: hello world
	// received output: this is a line
	// received output: and another line!
}

func ExampleStart_errWithStderr() {
	cmd := exec.Command("bash", "-o", "pipefail", "-o", "errexit", "-c",
		`echo "my stdout output" ; sleep 0.001 ; >&2 echo "my stderr output" ; exit 1`)
	stream, err := streamexec.Start(cmd, streamexec.Stdout|streamexec.ErrWithStderr)
	if err != nil {
		fmt.Println("failed to start: ", err.Error())
	}

	out, err := stream.String()
	fmt.Println("got error:", err.Error()) // error includes stderr
	fmt.Println("got output:", out)
	// Output:
	// got error: exit status 1: my stderr output
	// got output: my stdout output
}
