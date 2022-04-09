package canvas

import "testing"

func TestParseCoord(t *testing.T) {
	x, y, err := ParseCoord("\"134,1634\"")
	if err != nil {
		t.Fatal(err)
	}
	if x != 134 {
		t.Fatalf("expected x to be 134, got %v", x)
	}
	if y != 1634 {
		t.Fatalf("expected y to be 1634, got %v", y)
	}
}
