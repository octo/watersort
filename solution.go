package watersort

import (
	"container/heap"
	"fmt"
	"log"
	"math/rand"
)

type Solution struct {
	State    *State
	Steps    []Step
	Distance int
}

func (s Solution) Clone() Solution {
	state := s.State.Clone()
	steps := make([]Step, len(s.Steps))
	copy(steps, s.Steps)

	return Solution{
		State:    &state,
		Steps:    steps,
		Distance: s.Distance,
	}
}

func (s Solution) PossibleSteps() []Step {
	destinationsByColor := make(map[Color][]int)
	for i, b := range s.State.Bottles {
		if b.FreeSlots() == 0 {
			continue
		}
		tc := b.TopColor()
		destinationsByColor[tc] = append(destinationsByColor[tc], i)
	}

	var ret []Step
	for srcIndex, src := range s.State.Bottles {
		tc := src.TopColor()
		if tc == Empty {
			continue
		}

		for _, dstIndex := range append(destinationsByColor[tc], destinationsByColor[Empty]...) {
			if srcIndex == dstIndex {
				continue
			}

			ret = append(ret, Step{From: srcIndex, To: dstIndex})
		}
	}

	rand.Shuffle(len(ret), func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})

	return ret
}

type Step struct {
	From, To int
}

type Heap struct {
	Solutions []Solution
}

func (h Heap) Len() int {
	return len(h.Solutions)
}

func (h Heap) Less(i, j int) bool {
	return h.Solutions[i].Distance < h.Solutions[j].Distance
}

func (h *Heap) Swap(i, j int) {
	h.Solutions[i], h.Solutions[j] = h.Solutions[j], h.Solutions[i]
}

func (h *Heap) Push(x any) {
	h.Solutions = append(h.Solutions, x.(Solution))
}

func (h *Heap) Pop() any {
	last := h.Len() - 1
	s := h.Solutions[last]
	h.Solutions = h.Solutions[:last]
	return s
}

func FindSolution(s State) ([]Step, error) {
	sol := Solution{
		State: &s,
	}

	h := &Heap{}
	heap.Init(h)
	heap.Push(h, sol)

	seen := make(map[uint32]bool)

	for len(h.Solutions) > 0 {
		base := heap.Pop(h).(Solution)

		for _, step := range base.PossibleSteps() {
			next := base.Clone()

			if err := next.State.Pour(step.From, step.To); err != nil {
				log.Printf("State.Pour(%d, %d): %v", step.From, step.To, err)
				continue
			}

			chk := next.State.Checksum()
			if seen[chk] {
				continue
			}

			next.Steps = append(next.Steps, step)

			minRequiredMoves := next.State.MinRequiredMoves()
			next.Distance = len(next.Steps) + minRequiredMoves
			// log.Printf("Distance: %2d + %2d = %2d", len(next.Steps), minRequiredMoves, next.Distance)
			if minRequiredMoves == 0 {
				log.Printf("Evaluated %d states to find solution", len(seen))
				return next.Steps, nil
			}

			seen[chk] = true
			heap.Push(h, next)
		}
	}
	return nil, fmt.Errorf("unable to find solution; evaluated %d states", len(seen))
}
