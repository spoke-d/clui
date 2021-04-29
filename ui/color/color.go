package color

import (
	"fmt"
	"math"
)

// Color represents an underlying color that we can always render.
type Color interface {
	// RGB returns a color with 3 components.
	RGB() RGB
}

// Hex defines a hexadecimal value.
type Hex struct {
	r, g, b uint8
}

func (h Hex) String() string {
	return fmt.Sprintf("#%02x%02x%02x", h.r, h.g, h.b)
}

// RGB represents the sRGB (standard RGB) values in the range 0-1
type RGB struct {
	R, G, B float64
}

// MakeRGB creates a new RGB
func MakeRGB(r, g, b float64) RGB {
	return RGB{R: r, G: g, B: b}
}

func (c RGB) RGB() RGB {
	return c
}

func (c RGB) Hex() Hex {
	return Hex{
		r: uint8(c.R * 255),
		g: uint8(c.G * 255),
		b: uint8(c.B * 255),
	}
}

func (c RGB) HSL() HSL {
	min := math.Min(math.Min(c.R, c.G), c.B)
	max := math.Max(math.Max(c.R, c.G), c.B)

	hsl := HSL{
		L: (min + max) / 2,
	}

	if min == max {
		hsl.S = 0
		hsl.H = 0
	} else {
		if hsl.L < 0.5 {
			hsl.S = (max - min) / (max + min)
		} else {
			hsl.S = (max - min) / (2.0 - max - min)
		}

		if max == c.R {
			hsl.H = (c.G - c.B) / (max - min)
		} else if max == c.G {
			hsl.H = 2.0 + (c.B-c.R)/(max-min)
		} else {
			hsl.H = 4.0 + (c.R-c.G)/(max-min)
		}

		hsl.H *= 60

		if hsl.H < 0 {
			hsl.H += 360
		}
	}
	return hsl
}

func (c RGB) HSV() HSV {
	min := math.Min(math.Min(c.R, c.G), c.B)
	max := math.Max(math.Max(c.R, c.G), c.B)
	delta := max - min

	hsv := HSV{
		V: max,
	}

	if max != 0.0 {
		hsv.S = delta / max
	}

	hsv.H = 0.0 // We use 0 instead of undefined as in wp.
	if min != max {
		if max == c.R {
			hsv.H = math.Mod((c.G-c.B)/delta, 6.0)
		}
		if max == c.G {
			hsv.H = (c.B-c.R)/delta + 2.0
		}
		if max == c.B {
			hsv.H = (c.R-c.G)/delta + 4.0
		}
		hsv.H *= 60.0
		if hsv.H < 0.0 {
			hsv.H += 360.0
		}
	}

	return hsv
}

// Distance attempts to calculate the distance between two colours.
func (c RGB) Distance(o RGB) float64 {
	square := func(v float64) float64 {
		return v * v
	}
	return math.Sqrt(square(c.R-o.R) + square(c.G-o.G) + square(c.B-o.B))
}

// HSV represents the color hue, saturation and brighteness (value).
type HSV struct {
	H, S, V float64
}

func (c HSV) RGB() RGB {
	var (
		hp = c.H / 60.0
		v  = c.V * c.S
		x  = v * (1.0 - math.Abs(math.Mod(hp, 2.0)-1.0))

		m       = c.V - v
		r, g, b = 0.0, 0.0, 0.0
	)

	if hp >= 0.0 && hp < 1.0 {
		r = v
		g = x
	} else if hp >= 1.0 && hp < 2.0 {
		r = x
		g = v
	} else if hp >= 2.0 && hp < 3.0 {
		g = v
		b = x
	} else if hp >= 3.0 && hp < 4.0 {
		g = x
		b = v
	} else if hp >= 4.0 && hp < 5.0 {
		r = x
		b = v
	} else if hp >= 5.0 && hp < 6.0 {
		r = v
		b = x
	}

	return RGB{
		R: math.Round((m + r) * 255),
		G: math.Round((m + g) * 255),
		B: math.Round((m + b) * 255),
	}
}

// HSL represents the color hue, saturation and luminance
type HSL struct {
	H, S, L float64
}

func (c HSL) RGB() RGB {
	if c.S == 0 {
		return RGB{R: c.L, G: c.L, B: c.L}
	}

	var (
		h          float64
		r, g, b    float64
		t1, t2     float64
		tr, tg, tb float64
	)

	h = c.H

	if c.L < 0.5 {
		t1 = c.L * (1.0 + c.S)
	} else {
		t1 = c.L + c.S - c.L*c.S
	}

	t2 = 2*c.L - t1
	h /= 360
	tr = h + 1.0/3.0
	tg = h
	tb = h - 1.0/3.0

	if tr < 0 {
		tr++
	} else if tr > 1 {
		tr--
	}
	if tg < 0 {
		tg++
	} else if tg > 1 {
		tg--
	}
	if tb < 0 {
		tb++
	} else if tb > 1 {
		tb--
	}

	// Red
	if 6*tr < 1 {
		r = t2 + (t1-t2)*6*tr
	} else if 2*tr < 1 {
		r = t1
	} else if 3*tr < 2 {
		r = t2 + (t1-t2)*(2.0/3.0-tr)*6
	} else {
		r = t2
	}

	// Green
	if 6*tg < 1 {
		g = t2 + (t1-t2)*6*tg
	} else if 2*tg < 1 {
		g = t1
	} else if 3*tg < 2 {
		g = t2 + (t1-t2)*(2.0/3.0-tg)*6
	} else {
		g = t2
	}

	// Blue
	if 6*tb < 1 {
		b = t2 + (t1-t2)*6*tb
	} else if 2*tb < 1 {
		b = t1
	} else if 3*tb < 2 {
		b = t2 + (t1-t2)*(2.0/3.0-tb)*6
	} else {
		b = t2
	}

	return RGB{
		R: math.Round(r * 255),
		G: math.Round(g * 255),
		B: math.Round(b * 255),
	}
}
