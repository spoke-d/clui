package style

import (
	"testing"

	"github.com/spoke-d/clui/ui/color"
)

func TestParseStyle(t *testing.T) {
	style, err := ParseStyle("color=rgb(255,0,0);")
	if err != nil {
		t.Fatalf("expected error to be nil, received %v", err)
	}
	if style.color != color.MakeRGB(255, 0, 0) {
		t.Errorf("expected color to be rgb(255,0,0), received %v", style.color)
	}
}
