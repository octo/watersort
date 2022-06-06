package watersort

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPour(t *testing.T) {
	cases := []struct {
		name string
		in   State
		step Step
		want State
	}{
		{
			name: "into empty bottle",
			step: Step{From: 0, To: 3},
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
			step: Step{From: 2, To: 0},
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
			step: Step{From: 2, To: 0},
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
			step: Step{From: 3, To: 0},
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

			if err := tc.in.Apply(tc.step); err != nil {
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

func TestTestdata(t *testing.T) {
	path := filepath.Join("solver", "testdata")
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, de := range dirEntries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".json") {
			continue
		}

		t.Run(de.Name(), func(t *testing.T) {
			f, err := os.Open(filepath.Join(path, de.Name()))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			state, err := LoadLevel(f)
			if err != nil {
				t.Fatalf("LoadLevel(): %v", err)
			}

			if err := state.sanityCheck(); err != nil {
				t.Errorf("sanityCheck(): %v", err)
			}
		})
	}
}
