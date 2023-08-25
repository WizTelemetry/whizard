package monitoringgateway

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDifference(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"a", "b", "e"}

	if diff := cmp.Diff(difference(a, b), []string{"c", "d"}); diff != "" {
		t.Fatal(diff)
	}
	if diff := cmp.Diff(difference(b, a), []string{"e"}); diff != "" {
		t.Fatal(diff)
	}
}
