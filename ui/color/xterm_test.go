package color

import (
	"reflect"
	"testing"
)

func TestXTerm256Size(t *testing.T) {
	palette := XTerm256()
	size := palette.Size()
	// Handle the fact that there are duplicates
	if size != 247 {
		t.Errorf("expected size to be 256, received: %d", size)
	}
}

func TestXTerm256Nearest(t *testing.T) {
	palette := XTerm256()
	rgb, labels := palette.Nearest(MakeRGB(0, 0, 0))
	if rgb != MakeRGB(0, 0, 0) {
		t.Errorf("expected RGB{0,0,0}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Black", "Grey0"}) {
		t.Errorf("expected 'black, grey0', received %v", labels)
	}
}

func TestXTerm256NearestPlum(t *testing.T) {
	palette := XTerm256()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 255))
	if rgb != MakeRGB(255, 175, 255) {
		t.Errorf("expected RGB{255, 175, 255}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Plum1"}) {
		t.Errorf("expected 'Plum1', received %v", labels)
	}
}

func TestXTerm256NearestMiss(t *testing.T) {
	palette := XTerm256()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 240))
	if rgb != MakeRGB(255, 175, 255) {
		t.Errorf("expected RGB{255, 175, 255}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Plum1"}) {
		t.Errorf("expected 'Plum1', received %v", labels)
	}
}
