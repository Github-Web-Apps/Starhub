package diff_test

import (
	"testing"

	"github.com/caarlos0/watchub/shared/diff"
	"github.com/stretchr/testify/assert"
)

func TestSimpleDiff(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "b", "d", "e", "f"}
	t.Run("['c']", testDiffFunc(a, b, []string{"c"}))
	t.Run("['d','e','f']", testDiffFunc(b, a, []string{"d", "e", "f"}))
}

func testDiffFunc(left, right, expected []string) func(*testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, expected, diff.Of(left, right))
	}
}
