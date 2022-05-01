## Water Sort Solver

This package implements a Water Sort Puzzle solver.

## Getting started

First, you need to create a JSON file representing the level.
You can find example files in `solver/testdata/`.

Then run the executable in the `solver/` directory, for example:

```
$ cd solver
solver$ go build
solver$ ./solver -input=level.json
```

## Example run

This uses the provided sample data, Water Sort Puzzle's infamous level 105:

```
$ time ./solver <testdata/level105.json
2022/05/01 21:15:12 Evaluated 20637 states to find solution
Step  1: pour  6 onto 14
Step  2: pour  6 onto 13
Step  3: pour  7 onto  6
Step  4: pour 10 onto  6
Step  5: pour  5 onto 13
Step  6: pour 10 onto 14
Step  7: pour 11 onto 10
Step  8: pour 11 onto 14
Step  9: pour  8 onto 11
Step 10: pour  7 onto  8
Step 11: pour  7 onto  5
Step 12: pour  7 onto 11
Step 13: pour  2 onto  7
Step 14: pour  8 onto  7
Step 15: pour  2 onto  8
Step 16: pour 12 onto  2
Step 17: pour  9 onto  2
Step 18: pour  9 onto 10
Step 19: pour 12 onto  9
Step 20: pour  8 onto 12
Step 21: pour 11 onto  8
Step 22: pour  1 onto 11
Step 23: pour  1 onto 13
Step 24: pour  9 onto  1
Step 25: pour 11 onto  9
Step 26: pour 12 onto 11
Step 27: pour  3 onto 12
Step 28: pour  4 onto 12
Step 29: pour  3 onto  7
Step 30: pour  4 onto 11
Step 31: pour  4 onto  3
Step 32: pour  4 onto 12
Step 33: pour  3 onto  4
Step 34: pour  5 onto  4
Step 35: pour  6 onto  3
Step 36: pour  5 onto 14
Step 37: pour  9 onto  5
Step 38: pour  6 onto 13
Step 39: pour 10 onto  6
Step 40: pour  2 onto 10
Step 41: pour  6 onto  2

real    0m0.158s
user    0m0.047s
sys     0m0.078s
```

## Algorithm

This solver implements the [A* search algorithm](https://en.wikipedia.org/wiki/A*_search_algorithm)
to efficiently find a solution.

Basically the algorithm extends Dijkstra's path finding algorithm with a
heuristic to judge the distance to a solved state. The heuristic used is the 
lower bound of remaining moves. This is done in two steps:

1.  For each bottle, the number of colors that are on top of a different color
    is counted. For instance, a bottle with "Red, Green, Blue" needs at least
    two moves to sort the bottle, i.e. move Green and Blue out of the bottle.
1.  The distribution of the bottom-most colors is considered. For example, if
    Red is at the bottom of three bottles, at least two moves are required.

Other solutions I've seen use depth-first search (DFS). Using A* is advantageous
because:

*   The solution found by A* is always optimal.
*   The search space is dramatically reduced.

## License

Licensed under the ISC license

## Author

Florian Forster &lt;ff at octo.it&gt;