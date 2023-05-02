package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	ints := []int{1, 2, 3, 4}
	ints = Map(ints, func(i int) int { return i * 10 })
	assert.EqualValues(t, []int{10, 20, 30, 40}, ints)

	sl := []string{"a ", "b", " c"}
	sl = Map(sl, func(s string) string { return "#" + strings.TrimSpace(s) })
	assert.EqualValues(t, []string{"#a", "#b", "#c"}, sl)
}
