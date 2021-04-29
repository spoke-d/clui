package color

import (
	"reflect"
	"testing"
)

func TestTrueColorSize(t *testing.T) {
	palette := TrueColor()
	size := palette.Size()
	// Handle the fact that there are duplicates
	if size != 16777216 {
		t.Errorf("expected size to be 16777216, received: %d", size)
	}
}

func TestTrueColorNearest(t *testing.T) {
	palette := TrueColor()
	rgb, labels := palette.Nearest(MakeRGB(0, 0, 0))
	if rgb != MakeRGB(0, 0, 0) {
		t.Errorf("expected RGB{0,0,0}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Black", "Grey0"}) {
		t.Errorf("expected 'black, grey0', received %v", labels)
	}
}

func TestTrueColorNearestPlum(t *testing.T) {
	palette := TrueColor()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 255))
	if rgb != MakeRGB(255, 175, 255) {
		t.Errorf("expected RGB{255, 175, 255}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Plum1"}) {
		t.Errorf("expected 'Plum1', received %v", labels)
	}
}

func TestTrueColorNearestMiss(t *testing.T) {
	palette := TrueColor()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 240))
	if rgb != MakeRGB(255, 175, 240) {
		t.Errorf("expected RGB{255, 175, 240}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Plum1"}) {
		t.Errorf("expected 'Plum1', received %v", labels)
	}
}
