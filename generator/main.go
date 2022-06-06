package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/octo/watersort"
)

var (
	num  = flag.Int("num", 10, "number of colors/bottles; does not include empty bottles")
	size = flag.Int("size", 4, "number of slots in each bottle")
)

func main() {
	flag.Parse()

	maxComplexity := 0
	for {
		s := watersort.RandomState(*num, *size)

		var complexity int
		_, err := s.Solve(watersort.ReportComplexity(&complexity))
		if errors.Is(err, watersort.ErrNoSolution) {
			fmt.Println("=== Unsolvable ===")
			json.NewEncoder(os.Stdout).Encode(s)
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		if maxComplexity < complexity {
			fmt.Printf("=== Complexity %d ===\n", complexity)
			json.NewEncoder(os.Stdout).Encode(s)
			maxComplexity = complexity
		}
	}
}
