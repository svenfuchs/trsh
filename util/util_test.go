package util_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/util"
	"testing"
)

func TestCombinations(t *testing.T) {
	strs := []string{"a", "b", "c"}
	actual := util.Combinations(strs)
	expect := []string{
		"a",
		"a b",
		"a b c",
		"a c",
		"a c b",
		"b",
		"b a",
		"b a c",
		"b c",
		"b c a",
		"c",
		"c a",
		"c a b",
		"c b",
		"c b a",
	}
	assert.Equal(t, expect, actual)
}
