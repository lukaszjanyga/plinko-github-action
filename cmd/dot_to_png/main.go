package main

import (
	"os"

	"github.com/raishey/plinko/pkg/renderers"
)

func main() {
	fileName := os.Args[1]

	err := renderers.DotFileToImg(fileName, ".plinko/plinko.png", "png")

	if err != nil {
		panic(err)
	}
}
