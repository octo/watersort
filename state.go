package watersort

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

type State struct {
	Bottles []Bottle
}

func LoadLevel(r io.Reader) (State, error) {
	var s State
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return State{}, err
	}

	if err := s.sanityCheck(); err != nil {
		return State{}, err
	}

	return s, nil
}

func RandomState(colorsNum, bottleSize int) State {
	colors := make([]Color, colorsNum*bottleSize)
	for i := 0; i < colorsNum; i++ {
		for j := 0; j < bottleSize; j++ {
			colors[i*bottleSize+j] = Color(i + 1)
		}
	}

	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	var s State
	for i := 0; i < colorsNum; i++ {
		var b Bottle

		for j := 0; j < bottleSize; j++ {
			b.Colors = append(b.Colors, colors[i*bottleSize+j])
		}

		s.Bottles = append(s.Bottles, b)
	}

	var empty Bottle
	for j := 0; j < bottleSize; j++ {
		empty.Colors = append(empty.Colors, Empty)
	}
	s.Bottles = append(s.Bottles, empty.Clone(), empty.Clone())

	return s
}

func (s State) Clone() State {
	bottles := make([]Bottle, len(s.Bottles))
	for i, b := range s.Bottles {
		bottles[i] = b.Clone()
	}
	return State{
		Bottles: bottles,
	}
}

func (s State) sanityCheck() error {
	var bottleSize int
	colorCounts := make(map[Color]int)

	for i, b := range s.Bottles {
		if bottleSize == 0 {
			bottleSize = len(b.Colors)
		}
		if len(b.Colors) != bottleSize {
			return fmt.Errorf("not all bottles have the same size: bottle %d has %d colors, want %d",
				i+1, len(b.Colors), bottleSize)
		}
		for i, c := range b.Colors {
			colorCounts[c] = colorCounts[c] + 1
			if i != 0 && c != Empty && b.Colors[i-1] == Empty {
				return fmt.Errorf("bottle %d: cannot stack color on top of empty", i+1)
			}
		}
	}

	if colorCounts[Empty] != 2*bottleSize {
		return fmt.Errorf("got %d empty slots, want %d", colorCounts[Empty], 2*bottleSize)
	}
	for c, n := range colorCounts {
		if c == Empty {
			continue
		}
		if n != bottleSize {
			return fmt.Errorf("color %v: got %d slots, want %d", c, n, bottleSize)
		}
	}

	return nil
}

func (s *State) Apply(step Step) error {
	return s.Bottles[step.From].PourOnto(&s.Bottles[step.To])
}

func (s State) minRequiredMoves() int {
	var (
		ret          int
		bottomColors = make(map[Color]int, len(s.Bottles))
	)
	for _, b := range s.Bottles {
		ret += b.MinRequiredMoves()

		bc := b.BottomColor()
		bottomColors[bc] = bottomColors[bc] + 1
	}

	for c, cnt := range bottomColors {
		if c == Empty {
			continue
		}
		ret += (cnt - 1)
	}

	return ret
}

func (s State) Solved() bool {
	return s.minRequiredMoves() == 0
}

func (s State) checksum() uint32 {
	var data []byte
	for _, b := range s.Bottles {
		for _, c := range b.Colors {
			data = append(data, byte(c))
		}
	}

	return crc32.ChecksumIEEE(data)
}

func (s State) BottleSize() int {
	return len(s.Bottles[0].Colors)
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Bottles)
}

func (s *State) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.Bottles)
}

func (s State) MarshalText() ([]byte, error) {
	var bottles [][]byte
	for _, bottle := range s.Bottles {
		b, err := bottle.MarshalText()
		if err != nil {
			return nil, err
		}
		bottles = append(bottles, b)
	}

	return bytes.Join(bottles, []byte("_")), nil
}

func (s *State) UnmarshalText(text []byte) error {
	s.Bottles = nil

	bottles := bytes.Split(text, []byte("_"))
	for _, btxt := range bottles {
		var b Bottle
		if err := b.UnmarshalText(btxt); err != nil {
			return err
		}
		s.Bottles = append(s.Bottles, b)
	}

	return nil
}

type Bottle struct {
	Colors []Color
}

func (b Bottle) Clone() Bottle {
	colors := make([]Color, len(b.Colors))
	copy(colors, b.Colors)
	return Bottle{
		Colors: colors,
	}
}

func (b Bottle) TopColor() Color {
	for i := len(b.Colors) - 1; i >= 0; i-- {
		if b.Colors[i] != Empty {
			return b.Colors[i]
		}
	}
	return Empty
}

func (b Bottle) BottomColor() Color {
	return b.Colors[0]
}

func (b Bottle) TopColorCount() int {
	topColor := b.TopColor()
	ret := 0
	for i := len(b.Colors) - 1; i >= 0; i-- {
		if b.Colors[i] == Empty {
			continue
		}
		if b.Colors[i] != topColor {
			break
		}
		ret++
	}
	return ret
}

func (b Bottle) FreeSlots() int {
	for i := len(b.Colors) - 1; i >= 0; i-- {
		if b.Colors[i] != Empty {
			return len(b.Colors) - (i + 1)
		}
	}
	return len(b.Colors)
}

func (b *Bottle) PourOnto(other *Bottle) error {
	c := b.TopColor()

	n, err := other.add(c, b.TopColorCount())
	if err != nil {
		return err
	}

	for i := len(b.Colors) - 1; i >= 0 && n > 0; i-- {
		if b.Colors[i] == Empty {
			continue
		}
		if b.Colors[i] != c {
			return fmt.Errorf("cannot pop color %v", c)
		}
		b.Colors[i] = Empty
		n--
	}

	return nil
}

func (b *Bottle) add(c Color, n int) (int, error) {
	free := b.FreeSlots()
	if free == 0 {
		return 0, fmt.Errorf("no space available")
	}

	if tc := b.TopColor(); tc != Empty && tc != c {
		return 0, fmt.Errorf("cannot pour color %v onto %v", c, tc)
	}

	ret := n
	if n > free {
		ret = free
	}

	for i := 0; i < ret; i++ {
		index := len(b.Colors) + i - free
		if b.Colors[index] != Empty {
			return 0, fmt.Errorf("assertion failed: b.Colors[%d] = %v, want Empty", index, b.Colors[index])
		}
		b.Colors[index] = c
	}

	return ret, nil
}

func (b *Bottle) MinRequiredMoves() int {
	ret := 0
	for i := 1; i < len(b.Colors); i++ {
		if b.Colors[i] != Empty && b.Colors[i] != b.Colors[i-1] {
			ret++
		}
	}

	return ret
}

func (b Bottle) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Colors)
}

func (b *Bottle) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.Colors)
}

func (b Bottle) MarshalText() ([]byte, error) {
	var colors []string
	for _, c := range b.Colors {
		colors = append(colors, strconv.Itoa(int(c)))
	}
	return []byte(strings.Join(colors, "-")), nil
}

func (b *Bottle) UnmarshalText(text []byte) error {
	b.Colors = nil

	for _, c := range bytes.Split(text, []byte("-")) {
		color, err := strconv.Atoi(string(c))
		if err != nil {
			return err
		}
		b.Colors = append(b.Colors, Color(color))
	}

	return nil
}

type Color int

const (
	Empty Color = iota
	Blue
	Brown
	DarkBlue
	DarkGreen
	Gray
	Green
	LightBlue
	LightGreen
	Orange
	Pink
	Purple
	Red
	Yellow
)

var nameByColor = map[Color]string{
	Empty:      "Empty",
	Blue:       "Blue",
	Brown:      "Brown",
	DarkBlue:   "DarkBlue",
	DarkGreen:  "DarkGreen",
	Gray:       "Gray",
	Green:      "Green",
	LightGreen: "LightGreen",
	LightBlue:  "LightBlue",
	Orange:     "Orange",
	Pink:       "Pink",
	Purple:     "Purple",
	Red:        "Red",
	Yellow:     "Yellow",
}

func (c Color) String() string {
	if name, ok := nameByColor[c]; ok {
		return name
	}
	return fmt.Sprintf("color#%d", c)
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

var colorByName = map[string]Color{
	"Empty":      Empty,
	"Blue":       Blue,
	"Brown":      Brown,
	"DarkBlue":   DarkBlue,
	"DarkGreen":  DarkGreen,
	"Gray":       Gray,
	"Green":      Green,
	"LightGreen": LightGreen,
	"LightBlue":  LightBlue,
	"Orange":     Orange,
	"Pink":       Pink,
	"Purple":     Purple,
	"Red":        Red,
	"Yellow":     Yellow,
}

func (c *Color) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	var ok bool
	*c, ok = colorByName[name]
	if ok {
		return nil
	}

	i, err := strconv.Atoi(strings.TrimPrefix(name, "color#"))
	if err != nil {
		return fmt.Errorf("%q is not a valid color: %w", name, err)
	}

	*c = Color(i)
	return nil
}
