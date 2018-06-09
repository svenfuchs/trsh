package jq_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/jq"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/pipe"
	// "github.com/svenfuchs/trsh/spec"
	// "io"
	// "path"
	// "runtime"
	"testing"
	// "strings"
)

// var _, file, _, _ = runtime.Caller(0)
// var s = spec.Load(nil, path.Join(path.Dir(file), "../test/manifest.json"))

// Accept

func accept(str string, s ...string) (ops.Op, bool) {
	if len(s) == 0 {
		s = append(s, "repository find id=1")
	}
	// TODO needs the spec
	api.Init()
	jq.Init()
	pipe.Init()

	// a := api.New()
	// a.Accept(s[0])
	o := ops.New(conf.New("/tmp/trvs.test/config.json"))
	o.Accept(s[0] + " | " + str)
	op := o.Last()
	ok, _ := op.Accept(str)
	return op, ok
}

func acceptBuilds(str string, s ...string) (ops.Op, bool) {
	return accept(str, "builds find repository.id=1")
}

// func TestAcceptEmpty(t *testing.T) {
// 	j, ok := accept("")
// 	assert.False(t, ok)
// 	assert.Equal(t, "", j.Str())
// }

func TestAcceptDot(t *testing.T) {
	j, ok := accept(".")
	assert.True(t, ok)
	assert.Equal(t, ".", j.Str())
}

func TestAcceptDotAndUnknownAttr(t *testing.T) {
	_, ok := accept(".wat")
	assert.False(t, ok)
	// assert.Equal(t, ".", j.Str())
	// assert.Equal(t, "", j.Str())
}

func TestAcceptDotAndPartOfAttr(t *testing.T) {
	j, ok := accept(".na")
	assert.True(t, ok)
	assert.Equal(t, ".na", j.Str())
}

func TestAcceptDotAndAttr(t *testing.T) {
	j, ok := accept(".name")
	assert.True(t, ok)
	assert.Equal(t, ".name", j.Str())
}

func TestAcceptDotAndAttrAndSpace(t *testing.T) {
	j, ok := acceptBuilds(".builds ")
	assert.True(t, ok)
	assert.Equal(t, ".builds ", j.Str())
}

func TestAcceptDotAndAttrAndSpaceAndBrackets(t *testing.T) {
	j, ok := acceptBuilds(".builds []")
	assert.True(t, ok)
	assert.Equal(t, ".builds []", j.Str())
}

func TestAcceptDotAndAttrAndSpaceAndBracketsAndSpace(t *testing.T) {
	j, ok := acceptBuilds(".builds [] ")
	assert.True(t, ok)
	assert.Equal(t, ".builds [] ", j.Str())
}

func TestAcceptDotAndAttrAndSpaceAndBracketsAndSpaceAndDot(t *testing.T) {
	j, ok := acceptBuilds(".builds [] .")
	assert.True(t, ok)
	assert.Equal(t, ".builds [] .", j.Str())
}

func TestAcceptDotAndAttrAndSpaceAndBracketsAndSpaceAndDotAndAttr(t *testing.T) {
	j, ok := acceptBuilds(".builds [] .number")
	assert.True(t, ok)
	assert.Equal(t, ".builds [] .number", j.Str())
}

func TestAcceptDotAndAttrAndSpaceAndBracketsAndSpaceAndDotAndUnknownAttr(t *testing.T) {
	_, ok := acceptBuilds(".builds [] .wat")
	assert.False(t, ok)
	// assert.Equal(t, ".builds [] .wat", j.Str())
	// assert.Equal(t, "", j.Str())
}

// Hint

func hint(str string) string {
	if j, ok := accept(str); ok {
		return j.Hint(str)
	}
	return ""
}

func TestHintDot(t *testing.T) {
	assert.Equal(t, "id", hint("."))
}

func TestHintUnknownAttr(t *testing.T) {
	assert.Equal(t, "", hint(".wat"))
}

func TestHintPartOfAttr(t *testing.T) {
	assert.Equal(t, "me", hint(".na"))
}

func TestHintAttr(t *testing.T) {
	assert.Equal(t, "", hint(".name"))
}

// Completions

func completions(str string) []string {
	if j, ok := accept(str); ok {
		return j.Completions(str)
	}
	return []string{}
}

func TestCompletionsDot(t *testing.T) {
	c := []string{
		"id",
		"name",
		"slug",
		"description",
		"github_id",
		"github_language",
		"active",
		"private",
		"owner",
		"default_branch",
		"starred",
		"managed_by_installation",
		"active_on_org",
		"current_build",
		"last_started_build",
		"next_build_number",
	}
	assert.Equal(t, c, completions("."))
}

func TestCompletionsDotAttr(t *testing.T) {
	c := []string{
		"name",
	}
	assert.Equal(t, c, completions(".na"))
}
