package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/octo/watersort"
)

var (
	input = flag.String("input", "", "file to read from")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixMicro())

	var in io.Reader = os.Stdin
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			log.Fatalf("Open(%q): %v", *input, err)
		}
		in = f
	}

	level, err := watersort.LoadLevel(in)
	if err != nil {
		log.Fatalln("watersort.LoadLevel():", err)
	}

	steps, err := watersort.FindSolution(level)
	if err != nil {
		log.Fatalln("watersort.FindSolution():", err)
	}

	for i, step := range steps {
		fmt.Printf("Step %2d: pour %2d onto %2d\n", i+1, step.From+1, step.To+1)
	}
}
