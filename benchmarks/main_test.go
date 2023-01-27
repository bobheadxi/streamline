package benchmarks

import (
	"flag"
	"os"
	"testing"
)

var profile = flag.Bool("profile", false, "profile")

func TestMain(t *testing.M) {
	flag.Parse()
	os.Exit(t.Run())
}
