package wol

import "testing"

func TestListenWOL(t *testing.T) {
	var err error
	if err = ListenWOL(7); err != nil {
		t.Error(err)
	}
}

func TestWOL(t *testing.T) {
	go func() {
		if err := ListenWOL(7); err != nil {
			t.Fatal(err)
		}
	}()
	if err := WOL("10:7B:44:F2:90:31"); err != nil {
		t.Fatal(err)
	}
}
