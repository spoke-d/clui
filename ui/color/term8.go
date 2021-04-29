package color

// Term8 is a 8 color palette.
func Term8() *IndexedPalette {
	p := NewIndexedPalette()
	p.add(
		index{rgb: MakeRGB(0, 0, 0), labels: []string{"Black"}},
		index{rgb: MakeRGB(255, 0, 0), labels: []string{"Red"}},
		index{rgb: MakeRGB(0, 255, 0), labels: []string{"Green"}},
		index{rgb: MakeRGB(255, 255, 0), labels: []string{"Yellow"}},
		index{rgb: MakeRGB(0, 0, 255), labels: []string{"Blue"}},
		index{rgb: MakeRGB(255, 0, 255), labels: []string{"Magenta"}},
		index{rgb: MakeRGB(0, 255, 255), labels: []string{"Cyan"}},
		index{rgb: MakeRGB(255, 255, 255), labels: []string{"White"}},
	)
	return p
}
