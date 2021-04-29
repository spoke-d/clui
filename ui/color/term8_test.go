package color

import (
	"reflect"
	"testing"
)

func TestTerm8Size(t *testing.T) {
	palette := Term8()
	size := palette.Size()
	// Handle the fact that there are duplicates
	if size != 8 {
		t.Errorf("expected size to be 8, received: %d", size)
	}
}

func TestTerm8Nearest(t *testing.T) {
	palette := Term8()
	rgb, labels := palette.Nearest(MakeRGB(0, 0, 0))
	if rgb != MakeRGB(0, 0, 0) {
		t.Errorf("expected RGB{0,0,0}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"Black"}) {
		t.Errorf("expected 'black', received %v", labels)
	}
}

func TestTerm8NearestWhite(t *testing.T) {
	palette := Term8()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 255))
	if rgb != MakeRGB(255, 255, 255) {
		t.Errorf("expected RGB{255, 255, 255}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"White"}) {
		t.Errorf("expected 'White', received %v", labels)
	}
}

func TestTerm8NearestMiss(t *testing.T) {
	palette := Term8()
	rgb, labels := palette.Nearest(MakeRGB(255, 175, 240))
	if rgb != MakeRGB(255, 255, 255) {
		t.Errorf("expected RGB{255, 255, 255}, received: %v", rgb)
	}
	if !reflect.DeepEqual(labels, []string{"White"}) {
		t.Errorf("expected 'White', received %v", labels)
	}
}
