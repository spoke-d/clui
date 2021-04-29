package style

import (
	"testing"

	"github.com/spoke-d/clui/ui/color"
)

func TestParseColorStyle(t *testing.T) {
	style, err := ParseStyle("color=rgb(255,0,0);color=hsv(359,1,1.0);")
	if err != nil {
		t.Fatalf("expected error to be nil, received %v", err)
	}
	if style.color != color.MakeRGB(255, 0, 4) {
		t.Errorf("expected color to be rgb(255,0,0), received %v", style.color)
	}
}

func TestParseUnderlineStyle(t *testing.T) {
	style, err := ParseStyle("underline=true;")
	if err != nil {
		t.Fatalf("expected error to be nil, received %v", err)
	}
	if style.underline != true {
		t.Errorf("expected underline to be true, received %v", style.underline)
	}
}
