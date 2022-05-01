package watersort

import (
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

const level105OptimalSolution = 41

func TestFindSolution(t *testing.T) {
	cases := []struct {
		name string
		in   State
		want int
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
			want: 31,
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

			if got := len(steps); got != tc.want {
				t.Errorf("solution has %d steps, want %d", got, tc.want)
			}
		})
	}
}

func BenchmarkFindSolution(b *testing.B) {
	if err := level105.sanityCheck(); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		steps, err := FindSolution(level105.Clone())
		if err != nil {
			b.Error(err)
		}

		if got, want := len(steps), level105OptimalSolution; got != want {
			b.Errorf("got solution with %d steps, want %d", got, want)
		}
	}
}
