package diff_test

import (
	"testing"

	"github.com/caarlos0/watchub/internal/diff"
	"github.com/stretchr/testify/assert"
)

func TestSimpleDiff(t *testing.T) {
	assert := assert.New(t)
	a := []string{"a", "b", "c"}
	b := []string{"a", "b", "d", "e", "f"}
	assert.Equal([]string{"c"}, diff.Of(a, b))
	assert.Equal([]string{"d", "e", "f"}, diff.Of(b, a))
}
