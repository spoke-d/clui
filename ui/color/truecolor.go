package color

func TrueColor() *TrueColorPalette {
	return &TrueColorPalette{
		labelLocator: XTerm256(),
	}
}
