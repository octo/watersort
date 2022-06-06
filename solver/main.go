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
	input            = flag.String("input", "", "file to read from")
	reportComplexity = flag.Bool("report_complexity", false, "print how many states were considered to find the solution")
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

	var complexity int
	steps, err := level.Solve(watersort.ReportComplexity(&complexity))
	if err != nil {
		log.Fatalln("watersort.FindSolution():", err)
	}

	for i, step := range steps {
		fmt.Printf("Step %2d: %v\n", i+1, step)
	}
	if *reportComplexity {
		fmt.Printf("Complexity: %d\n", complexity)
	}
}
