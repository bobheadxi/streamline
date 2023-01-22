package benchmarks

import (
	"flag"
	"testing"
)

var profile = flag.Bool("profile", false, "profile")

func TestMain(t *testing.M) {
	flag.Parse()
}
