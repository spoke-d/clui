package color

import (
	"math"
	"sort"
)

// Palette defines a color palette for a series of colors.
type Palette struct {
	indexed map[RGB]struct{}
	colors  []RGB
}

// NewPalette creates a new indexed palette.
func NewPalette() *Palette {
	return &Palette{
		indexed: make(map[RGB]struct{}),
	}
}

// Add a RGB color to the index palette.
func (p *Palette) Add(rgb RGB) {
	p.add(rgb)
}

// Remove a RGB color from the index palette.
func (p *Palette) Remove(rgb RGB) {
	_, ok := p.indexed[rgb]
	if !ok {
		return
	}

	delete(p.indexed, rgb)

	for i, v := range p.colors {
		if v == rgb {
			p.colors = append(p.colors[:i], p.colors[i+1:]...)
			break
		}
	}
}

// Nearest finds the closest RGB color within the palette. Allows for converting
// between different color palettes.
func (p *Palette) Nearest(rgb RGB) RGB {
	var (
		result   RGB
		distance = math.MaxFloat64
	)
	for _, v := range p.colors {
		if d := v.Distance(rgb); d < distance {
			result = v
			distance = d
		}
	}
	return result
}

// Size returns the size of the palette.
func (p *Palette) Size() int {
	return len(p.colors)
}

// add is a bulk method for adding a lot of rgb colors with a better sorting.
func (p *Palette) add(colors ...RGB) {
	num := len(p.colors)
	for _, rgb := range colors {
		if _, ok := p.indexed[rgb]; ok {
			continue
		}

		p.indexed[rgb] = struct{}{}
		p.colors = append(p.colors, rgb)
	}
	if num != len(p.colors) {
		return
	}

	// Always sort the colors to get predictable results.
	sort.Slice(p.colors, func(i, j int) bool {
		hsl1 := p.colors[i].HSL()
		hsl2 := p.colors[j].HSL()
		if hsl1.H == hsl2.H {
			return hsl1.L < hsl2.L
		}
		return hsl1.H < hsl2.H
	})
}
