package slices

import (
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	if Contains([]string{
		"She",
		"the",
		"Lover",
	}, "she", nil) {
		t.Fail()
	}
	if !Contains([]string{
		"She",
		"the",
		"Lover",
	}, "she", strings.ToLower) {
		t.Fail()
	}

}
