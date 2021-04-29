package color

import (
	"math"
	"sort"
)

// Palette defines a palette of colors.
type Palette interface {
	// Add a color with associated labels. Multiple adds with the same color
	// will sum the labels together.
	Add(RGB, ...string)
	// Remove the color from the palette.
	Remove(RGB)
	// Nearest finds the nearest color associated with the color.
	Nearest(RGB) (RGB, []string)
	// Size returns the number of unique items within the palette.
	Size() int
}

type index struct {
	rgb    RGB
	labels []string
}

// IndexedPalette defines a color palette for a series of colors.
type IndexedPalette struct {
	indexed map[RGB]struct{}
	colors  []index
}

// NewIndexedPalette creates a new indexed palette.
func NewIndexedPalette() *IndexedPalette {
	return &IndexedPalette{
		indexed: make(map[RGB]struct{}),
	}
}

// Add a RGB color to the index palette.
func (p *IndexedPalette) Add(rgb RGB, labels []string) {
	p.add(index{
		rgb:    rgb,
		labels: labels,
	})
}

// Remove a RGB color from the index palette.
func (p *IndexedPalette) Remove(rgb RGB) {
	_, ok := p.indexed[rgb]
	if !ok {
		return
	}

	delete(p.indexed, rgb)

	for i, v := range p.colors {
		if v.rgb == rgb {
			p.colors = append(p.colors[:i], p.colors[i+1:]...)
			break
		}
	}
}

// Nearest finds the closest RGB color within the palette. Allows for converting
// between different color palettes.
func (p *IndexedPalette) Nearest(rgb RGB) (RGB, []string) {
	var (
		result   index
		distance = math.MaxFloat64
	)
	for _, v := range p.colors {
		if d := v.rgb.Distance(rgb); d < distance {
			result = v
			distance = d
		}
	}
	return result.rgb, result.labels
}

// Size returns the size of the palette.
func (p *IndexedPalette) Size() int {
	return len(p.colors)
}

// add is a bulk method for adding a lot of rgb colors with a better sorting.
func (p *IndexedPalette) add(indices ...index) {
	num := len(p.colors)
	for _, index := range indices {
		if _, ok := p.indexed[index.rgb]; ok {
			i := indexOf(p.colors, index.rgb)
			if i == -1 {
				panic("cache missmatch")
			}
			p.colors[i].labels = union(p.colors[i].labels, index.labels)
			continue
		}

		p.indexed[index.rgb] = struct{}{}
		p.colors = append(p.colors, index)
	}
	if num != len(p.colors) {
		return
	}

	// Always sort the colors to get predictable results.
	sort.Slice(p.colors, func(i, j int) bool {
		hsl1 := p.colors[i].rgb.HSL()
		hsl2 := p.colors[j].rgb.HSL()
		if hsl1.H == hsl2.H {
			return hsl1.L < hsl2.L
		}
		return hsl1.H < hsl2.H
	})
}

func indexOf(haystack []index, needle RGB) int {
	for i, v := range haystack {
		if v.rgb == needle {
			return i
		}
	}
	return -1
}

func union(a, b []string) []string {
	x := make(map[string]struct{})
	for _, v := range a {
		x[v] = struct{}{}
	}
	for _, v := range b {
		x[v] = struct{}{}
	}
	y := make([]string, 0, len(x))
	for v := range x {
		y = append(y, v)
	}
	sort.Strings(y)
	return y
}

//

// TrueColorPalette defines a color palette for a series of colors.
type TrueColorPalette struct{}

// NewTrueColorPalette creates a new indexed palette.
func NewTrueColorPalette() *TrueColorPalette {
	return &TrueColorPalette{}
}

// Add a RGB color to the index palette.
func (p *TrueColorPalette) Add(rgb RGB, labels []string) {}

// Remove a RGB color from the index palette.
func (p *TrueColorPalette) Remove(rgb RGB) {}

// Nearest finds the closest RGB color within the palette. Allows for converting
// between different color palettes.
func (p *TrueColorPalette) Nearest(rgb RGB) (RGB, []string) {
	return rgb, nil
}

// Size returns the size of the palette.
func (p *TrueColorPalette) Size() int {
	return 16777216
}
