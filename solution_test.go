package watersort

import (
	"log"
	"testing"
)

// level105 is "Water Sort Puzzle"'s infamous 105th level.
var level105 = State{
	Bottles: []Bottle{
		{Colors: []Color{Blue, Blue, DarkGreen, DarkBlue}},
		{Colors: []Color{Gray, Green, Pink, Purple}},
		{Colors: []Color{Brown, Red, Purple, Orange}},
		{Colors: []Color{Orange, Red, Pink, Orange}},
		{Colors: []Color{DarkBlue, Yellow, Red, DarkGreen}},
		{Colors: []Color{DarkGreen, Brown, DarkGreen, Yellow}},
		{Colors: []Color{LightGreen, Red, Purple, Brown}},
		{Colors: []Color{LightGreen, Pink, Purple, LightGreen}},
		{Colors: []Color{DarkBlue, Blue, Gray, Green}},
		{Colors: []Color{Green, Gray, Yellow, Brown}},
		{Colors: []Color{DarkBlue, LightGreen, Yellow, Gray}},
		{Colors: []Color{Orange, Pink, Blue, Green}},
		{Colors: []Color{Empty, Empty, Empty, Empty}},
		{Colors: []Color{Empty, Empty, Empty, Empty}},
	},
}

func TestFindSolution(t *testing.T) {
	cases := []struct {
		name string
		in   State
		// want []Step
	}{
		{
			name: "solve puzzle",
			in: State{
				Bottles: []Bottle{
					{Colors: []Color{Pink, Yellow, Purple, Orange}},
					{Colors: []Color{DarkGreen, Pink, Blue, Red}},
					{Colors: []Color{DarkGreen, DarkBlue, DarkBlue, Red}},
					{Colors: []Color{DarkBlue, Gray, Pink, Gray}},
					{Colors: []Color{Blue, Purple, Blue, Purple}},
					{Colors: []Color{Green, Red, DarkBlue, Orange}},
					{Colors: []Color{Yellow, DarkGreen, Orange, Gray}},
					{Colors: []Color{Orange, Green, Green, Gray}},
					{Colors: []Color{Red, Yellow, DarkGreen, Blue}},
					{Colors: []Color{Pink, Yellow, Green, Purple}},
					{Colors: []Color{Empty, Empty, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty, Empty}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.in.sanityCheck(); err != nil {
				t.Fatal(err)
			}

			steps, err := FindSolution(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("Number of steps: %d", len(steps))
			for i, step := range steps {
				log.Printf("Step %2d: pour %2d onto %2d", i+1, step.From+1, step.To+1)
			}
		})
	}
}
