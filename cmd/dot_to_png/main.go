package main

import (
	"os"

	"github.com/raishey/plinko/pkg/renderers"
)

func main() {
	fileName := os.Args[1]

	renderers.DotFileToImg(fileName, ".plinko/plinko.png", "png")
}
