package wol

import "testing"

func TestWOL(t *testing.T) {
	go func() {
		if err := ListenWOL(7); err != nil {
			t.Fatal(err)
		}
	}()
	if err := WOL("12345678ABCD"); err != nil {
		t.Fatal(err)
	}
}
