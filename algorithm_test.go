package mselector

import "testing"

func TestIsLack(t *testing.T) {
	if !isLack(1, 2) {
		t.Error("is lack error")
	}
	if !isLack(11, 12) {
		t.Error("is lack error")
	}
	if !isLack(21, 22) {
		t.Error("is lack error")
	}
	if isLack(11, 21) {
		t.Error("is lack error")
	}
	if isLack(21, 0) {
		t.Error("is lack error")
	}
}
