package main

import (
	"fmt"

	buildOpts "github.com/water-sucks/optnix/internal/build"
)

func main() {
	fmt.Printf("optnix version %v\n", buildOpts.Version)
}
