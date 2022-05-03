package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/octo/watersort"
)

var (
	num  = flag.Int("num", 10, "number of colors/bottles; does not include empty bottles")
	size = flag.Int("size", 4, "number of slots in each bottle")
)

func randomState(num, size int) watersort.State {
	flag.Parse()

	rand.Seed(time.Now().UnixMicro())

	colors := make([]watersort.Color, num*size)
	for i := 0; i < num; i++ {
		for j := 0; j < size; j++ {
			colors[i*size+j] = watersort.Color(i + 1)
		}
	}

	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	var s watersort.State
	for i := 0; i < num; i++ {
		var b watersort.Bottle

		for j := 0; j < size; j++ {
			b.Colors = append(b.Colors, colors[i*size+j])
		}

		s.Bottles = append(s.Bottles, b)
	}

	var empty watersort.Bottle
	for j := 0; j < size; j++ {
		empty.Colors = append(empty.Colors, watersort.Empty)
	}
	s.Bottles = append(s.Bottles, empty.Clone(), empty.Clone())

	return s
}

func main() {
	var maxComplexity int

	for {
		s := randomState(*num, *size)

		var complexity int
		_, err := watersort.FindSolution(s, watersort.ReportComplexity(&complexity))
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
