package hash_test

import (
	"ZakkBob/AskDave/crawler/hash"
	"testing"
)

func TestHashs(t *testing.T) {
	s := "Testing..."
	hash := hash.Hashs(s)
	got := hash.String()
	expected := "9d770c909c2c69b09eae2372c4cf405d"
	if got != expected {
		t.Errorf("expected '%s', but got '%s'", expected, got)
	}
}
