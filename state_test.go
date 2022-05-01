package watersort

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPour(t *testing.T) {
	cases := []struct {
		name     string
		in       State
		from, to int
		want     State
	}{
		{
			name: "into empty bottle",
			from: 0,
			to:   3,
			in: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Blue}},
					{Colors: []Color{Green, Blue, Red}},
					{Colors: []Color{Blue, Red, Green}},
					{Colors: []Color{Empty, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
			want: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Empty}},
					{Colors: []Color{Green, Blue, Red}},
					{Colors: []Color{Blue, Red, Green}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
		},
		{
			name: "into non-empty bottle",
			from: 2,
			to:   0,
			in: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Empty}},
					{Colors: []Color{Green, Blue, Red}},
					{Colors: []Color{Blue, Red, Green}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
			want: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Green}},
					{Colors: []Color{Green, Blue, Red}},
					{Colors: []Color{Blue, Red, Empty}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
		},
		{
			name: "partial pour",
			from: 2,
			to:   0,
			in: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Empty}},
					{Colors: []Color{Red, Blue, Red}},
					{Colors: []Color{Blue, Green, Green}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
			want: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Green, Green}},
					{Colors: []Color{Red, Blue, Red}},
					{Colors: []Color{Blue, Green, Empty}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
		},
		{
			name: "clean out a bottle",
			from: 3,
			to:   0,
			in: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Blue, Empty}},
					{Colors: []Color{Red, Green, Red}},
					{Colors: []Color{Blue, Green, Green}},
					{Colors: []Color{Blue, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
			want: State{
				Bottles: []Bottle{
					{Colors: []Color{Red, Blue, Blue}},
					{Colors: []Color{Red, Green, Red}},
					{Colors: []Color{Blue, Green, Green}},
					{Colors: []Color{Empty, Empty, Empty}},
					{Colors: []Color{Empty, Empty, Empty}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.in.sanityCheck(); err != nil {
				t.Error(err)
			}

			if err := tc.in.Pour(tc.from, tc.to); err != nil {
				t.Fatal(err)
			}

			if err := tc.in.sanityCheck(); err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(tc.want, tc.in); diff != "" {
				t.Errorf("state differs (-want/+got):\n%s", diff)
			}
		})
	}
}

func TestBottle_MinRequiredMoves(t *testing.T) {
	cases := []struct {
		name   string
		colors []Color
		want   int
	}{
		{
			name:   "max",
			colors: []Color{Red, Green, Blue, Red},
			want:   3,
		},
		{
			name:   "duplicate color",
			colors: []Color{Red, Green, Green, Blue},
			want:   2,
		},
		{
			name:   "done stack",
			colors: []Color{Red, Red, Red, Red},
			want:   0,
		},
		{
			name:   "mostly done stack",
			colors: []Color{Red, Red, Red, Empty},
			want:   0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Bottle{
				Colors: tc.colors,
			}

			if got := b.MinRequiredMoves(); got != tc.want {
				t.Errorf("MinRequiredMoves() = %d, want %d", got, tc.want)
			}
		})
	}
}
