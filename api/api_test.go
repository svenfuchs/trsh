package api_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/ops"
	s "github.com/svenfuchs/trsh/spec"
	p "path"
	"runtime"
	"testing"
)

var _, file, _, _ = runtime.Caller(0)
var path = p.Join(p.Dir(file), "../test/manifest.json")
var spec = s.Load(nil, path)

// Accept

func accept(str string) (ops.Op, bool, string) {
	api.Init(spec)
	a := api.New()
	ok, str := a.Accept(str)
	return a, ok, str
}

func TestAcceptEmpty(t *testing.T) {
	a, ok, str := accept("")
	assert.False(t, ok)
	assert.Equal(t, "", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptUnknownResource(t *testing.T) {
	a, ok, str := accept("wat")
	assert.False(t, ok)
	assert.Equal(t, "", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfResource(t *testing.T) {
	a, ok, str := accept("req")
	assert.True(t, ok)
	assert.Equal(t, "request", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfResourceAndSpace(t *testing.T) {
	a, ok, str := accept("req ")
	assert.False(t, ok)
	assert.Equal(t, "request", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptResource(t *testing.T) {
	a, ok, str := accept("user")
	assert.True(t, ok)
	assert.Equal(t, "user", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptResourceAndSpace(t *testing.T) {
	a, ok, str := accept("user ")
	assert.True(t, ok)
	assert.Equal(t, "user", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptUnknownAction(t *testing.T) {
	a, ok, str := accept("user wat")
	assert.False(t, ok)
	assert.Equal(t, "user", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfAction(t *testing.T) {
	a, ok, str := accept("user fi")
	assert.True(t, ok)
	assert.Equal(t, "user", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfActionAndSpace(t *testing.T) {
	a, ok, str := accept("user fi ")
	assert.False(t, ok)
	assert.Equal(t, "user", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptAction(t *testing.T) {
	a, ok, str := accept("user find")
	assert.True(t, ok)
	assert.Equal(t, "user find", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptActionAndSpace(t *testing.T) {
	a, ok, str := accept("user find ")
	assert.True(t, ok)
	assert.Equal(t, "user find", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptUnknownParamName(t *testing.T) {
	a, ok, str := accept("user find wat")
	assert.False(t, ok)
	assert.Equal(t, "user find", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfParamName(t *testing.T) {
	a, ok, str := accept("user find i")
	assert.True(t, ok)
	assert.Equal(t, "user find id=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptPartOfParamNameAndSpace(t *testing.T) {
	a, ok, str := accept("user find i ")
	assert.False(t, ok)
	assert.Equal(t, "user find", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamName(t *testing.T) {
	a, ok, str := accept("user find id")
	assert.True(t, ok)
	assert.Equal(t, "user find id=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamNameAndEqual(t *testing.T) {
	a, ok, str := accept("user find id=")
	assert.True(t, ok)
	assert.Equal(t, "user find id=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParam(t *testing.T) {
	a, ok, str := accept("user find id=1")
	assert.True(t, ok)
	assert.Equal(t, "user find id=1", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamAndSpace(t *testing.T) {
	a, ok, str := accept("user find id=1 ")
	assert.True(t, ok)
	assert.Equal(t, "user find id=1", a.Str())
	assert.Equal(t, " ", str)
}

func TestAcceptParamAndSpaceAndPipe(t *testing.T) {
	a, ok, str := accept("user find id=1 |")
	assert.True(t, ok)
	assert.Equal(t, "user find id=1", a.Str())
	assert.Equal(t, " |", str)
}

func TestAcceptParamAndSpaceAndUnknownParamName(t *testing.T) {
	a, ok, str := accept("builds find repository.id=1 wat")
	assert.False(t, ok)
	assert.Equal(t, "builds find repository.id=1", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamAndSpaceAndPartOfParamName(t *testing.T) {
	a, ok, str := accept("builds find repository.id=1 lim")
	assert.True(t, ok)
	assert.Equal(t, "builds find repository.id=1 limit=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamAndSpaceAndParamName(t *testing.T) {
	a, ok, str := accept("builds find repository.id=1 limit")
	assert.True(t, ok)
	assert.Equal(t, "builds find repository.id=1 limit=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamAndSpaceAndParamNameAndEqual(t *testing.T) {
	a, ok, str := accept("builds find repository.id=1 limit=")
	assert.True(t, ok)
	assert.Equal(t, "builds find repository.id=1 limit=", a.Str())
	assert.Equal(t, "", str)
}

func TestAcceptParamAndSpaceAndParam(t *testing.T) {
	a, ok, str := accept("builds find repository.id=1 limit=20")
	assert.True(t, ok)
	assert.Equal(t, "builds find repository.id=1 limit=20", a.Str())
	assert.Equal(t, "", str)
}

// Complete

func complete(str string) bool {
	api.Init(spec)
	a := api.New()
	a.Accept(str)
	return a.Complete()
}

func TestCompleteEmpty(t *testing.T) {
	assert.False(t, complete(""))
}

func TestCompletePartOfResource(t *testing.T) {
	assert.False(t, complete("req"))
}

func TestCompleteResource(t *testing.T) {
	assert.False(t, complete("repository"))
}

func TestCompleteResourceAndSpace(t *testing.T) {
	assert.False(t, complete("repository "))
}

func TestCompletePartOfAction(t *testing.T) {
	assert.False(t, complete("repository fi"))
}

func TestCompleteAction(t *testing.T) {
	assert.False(t, complete("repository find"))
}

func TestCompleteActionAndSpace(t *testing.T) {
	assert.False(t, complete("repository find "))
}

func TestCompletePartOfParamName(t *testing.T) {
	assert.False(t, complete("repository find i"))
}

func TestCompleteParamName(t *testing.T) {
	assert.False(t, complete("repository find id"))
}

func TestCompleteParamNameAndEquals(t *testing.T) {
	assert.False(t, complete("repository find id="))
}

func TestCompleteParam(t *testing.T) {
	assert.True(t, complete("repository find id=1"))
}

func TestCompleteParamAndSpace(t *testing.T) {
	assert.True(t, complete("builds find repository.id=1 "))
}

func TestCompleteParamAndPartOfParamName(t *testing.T) {
	assert.False(t, complete("builds find repository.id=1 l"))
}

func TestCompleteParamAndParamName(t *testing.T) {
	assert.False(t, complete("builds find repository.id=1 limit"))
}

func TestCompleteParamAndParamNameAndEqual(t *testing.T) {
	assert.False(t, complete("builds find repository.id=1 limit="))
}

func TestCompleteParamAndParamNameAndEqualAndValue(t *testing.T) {
	assert.True(t, complete("builds find repository.id=1 limit=1"))
}

// Hint

func hint(str string) string {
	api.Init(spec)
	a := api.New()
	a.Accept(str)
	return a.Hint(str)
}

func TestHintEmpty(t *testing.T) {
	h := hint("")
	assert.Equal(t, "active for_owner organization.login=", h)
}

func TestHintPartOfResource(t *testing.T) {
	h := hint("req")
	assert.Equal(t, "uest find id=", h)
}

func TestHintResource(t *testing.T) {
	h := hint("repository")
	assert.Equal(t, " find id=", h)
}

func TestHintResourceAndSpace(t *testing.T) {
	h := hint("repository ")
	assert.Equal(t, "find id=", h)
}

func TestHintPartOfAction(t *testing.T) {
	h := hint("repository fi")
	assert.Equal(t, "nd id=", h)
}

func TestHintAction(t *testing.T) {
	h := hint("repository find")
	assert.Equal(t, " id=", h)
}

func TestHintActionAndSpace(t *testing.T) {
	h := hint("repository find ")
	assert.Equal(t, "id=", h)
}

func TestHintPartOfParamName(t *testing.T) {
	h := hint("repository find i")
	assert.Equal(t, "d=", h)
}

func TestHintParamName(t *testing.T) {
	h := hint("repository find id")
	assert.Equal(t, "=", h)
}

func TestHintParamNameAndEquals(t *testing.T) {
	h := hint("repository find id=")
	assert.Equal(t, "", h)
}

func TestHintParam(t *testing.T) {
	h := hint("repository find id=1")
	assert.Equal(t, "", h)
}

func TestHintParamAndSpace(t *testing.T) {
	h := hint("builds find repository.id=1 ")
	assert.Equal(t, "branch.name=", h)
}

func TestHintParamAndPartOfParamName(t *testing.T) {
	h := hint("builds find repository.id=1 l")
	assert.Equal(t, "imit=", h)
}

func TestHintParamAndParamName(t *testing.T) {
	h := hint("builds find repository.id=1 limit")
	assert.Equal(t, "=", h)
}

func TestHintParamAndParamNameAndEqual(t *testing.T) {
	h := hint("builds find repository.id=1 limit=")
	assert.Equal(t, "", h)
}

func TestHintParamAndParamNameAndEqualAndValue(t *testing.T) {
	h := hint("builds find repository.id=1 limit=1")
	assert.Equal(t, "", h)
}

// Completions

func completions(str string) []string {
	api.Init(spec)
	a := api.New()
	a.Accept(str)
	return a.Completions(str)
}

func TestApiCompletionsEmpty(t *testing.T) {
	c := []string{
		"active",
		"beta_feature",
		"beta_features",
		"branch",
		"branches",
		"broadcast",
		"broadcasts",
		"build",
		"builds",
		"caches",
		"commit",
		"cron",
		"crons",
		"env_var",
		"env_vars",
		"error",
		"home",
		"installation",
		"job",
		"jobs",
		"key_pair",
		"key_pair_generated",
		"lint",
		"log",
		"message",
		"messages",
		"organization",
		"organizations",
		"owner",
		"repositories",
		"repository",
		"request",
		"requests",
		"resource",
		"setting",
		"settings",
		"stage",
		"stages",
		"template",
		"user",
	}
	assert.Equal(t, c, completions(""))
}

func TestApiCompletionsPartOfResource(t *testing.T) {
	c := []string{
		"repositories",
		"repository",
	}
	assert.Equal(t, c, completions("rep"))
}

func TestApiCompletionsResource(t *testing.T) {
	c := []string{
		"repository ",
	}
	assert.Equal(t, c, completions("repository"))
}

func TestApiCompletionsResourceAndSpace(t *testing.T) {
	c := []string{
		"find",
		"activate",
		"deactivate",
		"star",
		"unstar",
	}
	assert.Equal(t, c, completions("repository "))
}

func TestApiCompletionsPartOfAction(t *testing.T) {
	c := []string{
		"find ",
	}
	assert.Equal(t, c, completions("repository fi"))
}

func TestApiCompletionsAction(t *testing.T) {
	c := []string{
		"find ",
	}
	assert.Equal(t, c, completions("repository find"))
}

func TestApiCompletionsActionAndSpace(t *testing.T) {
	c := []string{
		"id=",
		"slug=",
	}
	assert.Equal(t, c, completions("repository find "))
}

func TestApiCompletionsPartOfParamName(t *testing.T) {
	c := []string{
		"id=",
	}
	assert.Equal(t, c, completions("repository find i"))
}

func TestApiCompletionsParamName(t *testing.T) {
	c := []string{
		"id=",
	}
	assert.Equal(t, c, completions("repository find id"))
}

func TestApiCompletionsParamNameAndEqual(t *testing.T) {
	c := []string{}
	assert.Equal(t, c, completions("repository find id="))
}

func TestApiCompletionsParam(t *testing.T) {
	c := []string{}
	assert.Equal(t, c, completions("repository find id=1"))
}

func TestApiCompletionsParamAndPartOfUnknownParamName(t *testing.T) {
	c := []string{}
	assert.Equal(t, c, completions("repository find id=1 r"))
}

func TestApiCompletionsParamAndPartOfParamName(t *testing.T) {
	c := []string{
		"limit=",
	}
	assert.Equal(t, c, completions("builds find repository.id=1 l"))
}

func TestApiCompletionsParamAndParamName(t *testing.T) {
	c := []string{
		"limit=",
	}
	assert.Equal(t, c, completions("builds find repository.id=1 limit"))
}

func TestApiCompletionsParamAndParamNameAndEqual(t *testing.T) {
	c := []string{}
	assert.Equal(t, c, completions("builds find repository.id=1 limit="))
}

func TestApiCompletionsParamAndParamNameAndEqualAndValue(t *testing.T) {
	c := []string{}
	assert.Equal(t, c, completions("builds find repository.id=1 limit=1"))
}
