package prompt_test

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/led-go"
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/jq"
	"github.com/svenfuchs/trsh/pipe"
	"github.com/svenfuchs/trsh/prompt"
	"strings"
	"testing"
	"time"
)

func TestPrompt(t *testing.T) {
	prompt, term := promptSetup()
	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos())
	assertOut(t, term,
		"<cr><clear>t ~ <cr><rgt-4>",
	)
}

// Keys

func TestPromptUnknownResource(t *testing.T) {
	prompt, term := promptSetup()
	receive(term, "foo")

	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos())
	assertOut(t, term,
		"<cr><clear>t ~ <cr><rgt-4>",
		"<red>foo<reset>",
	)
}

// Hints

func TestPromptHints(t *testing.T) {
	prompt, term := promptSetup()

	receive(term, "us")
	assert.Equal(t, 2, prompt.Pos())
	assert.Equal(t, "us", prompt.Str())
	assert.Equal(t, "er find id=", prompt.Suggested())

	receive(term, "er")
	assert.Equal(t, "user", prompt.Str())
	assert.Equal(t, " find id=", prompt.Suggested())

	receive(term, " ")
	assert.Equal(t, "user ", prompt.Str())
	assert.Equal(t, "find id=", prompt.Suggested())

	receive(term, "fi")
	assert.Equal(t, "user fi", prompt.Str())
	assert.Equal(t, "nd id=", prompt.Suggested())

	receive(term, "nd")
	assert.Equal(t, "user find", prompt.Str())
	assert.Equal(t, " id=", prompt.Suggested())

	receive(term, " ")
	assert.Equal(t, "user find ", prompt.Str())
	assert.Equal(t, "id=", prompt.Suggested())

	receive(term, "i")
	assert.Equal(t, "user find i", prompt.Str())
	assert.Equal(t, "d=", prompt.Suggested())

	receive(term, "d")
	assert.Equal(t, "user find id", prompt.Str())
	assert.Equal(t, "=", prompt.Suggested())

	receive(term, "=")
	assert.Equal(t, "user find id=", prompt.Str())
	assert.Equal(t, "", prompt.Suggested())

	receive(term, "1")
	assert.Equal(t, "user find id=1", prompt.Str())
	assert.Equal(t, "", prompt.Suggested())

	receive(term, " ")
	assert.Equal(t, "user find id=1 ", prompt.Str())
	assert.Equal(t, "|", prompt.Suggested())

	receive(term, "|")
	assert.Equal(t, "user find id=1 |", prompt.Str())
	assert.Equal(t, "", prompt.Suggested())

	receive(term, " ")
	assert.Equal(t, "user find id=1 | ", prompt.Str())
	assert.Equal(t, "", prompt.Suggested())

	receive(term, ".")
	assert.Equal(t, "user find id=1 | .", prompt.Str())
	assert.Equal(t, "id", prompt.Suggested())

	receive(term, "l")
	assert.Equal(t, "user find id=1 | .l", prompt.Str())
	assert.Equal(t, "ogin", prompt.Suggested())
}

// Completion

func TestPromptCompletion(t *testing.T) {
	prompt, term := promptSetup()

	receive(term, key(led.Tab))
	assert.Equal(t, 6, prompt.Pos())
	assert.Equal(t, "active", prompt.Str())
	times(6, func(_ int) { receive(term, key(led.Backspace)) })

	receive(term, "re")
	receive(term, key(led.Tab))
	assert.Equal(t, 12, prompt.Pos())
	assert.Equal(t, "repositories", prompt.Str())
	assert.Equal(t, " for_current_user active=", prompt.Suggested())

	receive(term, key(led.Tab))
	assert.Equal(t, 10, prompt.Pos())
	assert.Equal(t, "repository", prompt.Str())
	assert.Equal(t, " find id=", prompt.Suggested())

	receive(term, key(led.Tab))
	assert.Equal(t, "request", prompt.Str())

	receive(term, key(led.ShiftTab))
	assert.Equal(t, "repository", prompt.Str())

	receive(term, " ")
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find", prompt.Str())

	receive(term, key(led.Tab))
	assert.Equal(t, "repository activate", prompt.Str())

	receive(term, key(led.ShiftTab))
	assert.Equal(t, "repository find", prompt.Str())

	receive(term, " ")
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find id=", prompt.Str())

	receive(term, key(led.Tab))
	assert.Equal(t, "repository find slug=", prompt.Str())
}

func TestPromptCompletionAfterBackspaceOnAction(t *testing.T) {
	prompt, term := promptSetup()

	receive(term, "re")
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	receive(term, " ")
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find", prompt.Str())

	times(4, func(_ int) { receive(term, key(led.Backspace)) })
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find", prompt.Str())

	times(5, func(_ int) { receive(term, key(led.Backspace)) })
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find", prompt.Str())

	times(6, func(_ int) { receive(term, key(led.Backspace)) })
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	assert.Equal(t, "repository", prompt.Str())
}

func TestCompletionAfterBackspaceOnParam(t *testing.T) {
	prompt, term := promptSetup()

	receive(term, "repository ")
	receive(term, "find")
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find id=", prompt.Str())

	receive(term, key(led.Backspace))
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	receive(term, key(led.Tab))
	assert.Equal(t, "repository find id=", prompt.Str())
}

func assertOut(t *testing.T, term *testTerm, strs ...string) {
	out := string(led.Deansi([]byte(term.out)))
	assert.Equal(t, strings.Join(strs, ""), out)
}

func promptSetup() (*prompt.Prompt, *testTerm) {
	// c := http.New()
	// s := spec.Load(c)
	// TODO needs the spec
	api.Init()
	jq.Init()
	pipe.Init()
	c := conf.New("/tmp/trsh.test/config.json")
	t := newTestTerm()
	p := prompt.New(c, "t ~ ", t)
	go p.Run()
	time.Sleep(1 * time.Millisecond)
	return p, t
}

func receive(t *testTerm, str string) {
	t.keys <- str
	time.Sleep(5 * time.Millisecond)
}

func reset(term *testTerm) {
	term.out = ""
}

func times(n int, f func(int)) {
	for i := 0; i < n; i++ {
		f(i)
	}
}

func p(s string) {
	println("\"" + s + "\"")
}

func key(k int) string {
	return led.Keys[k].Str()
}

func newTestTerm() *testTerm {
	k := make(chan string, 1)
	t := testTerm{keys: k}
	return &t
}

type testTerm struct {
	keys chan string
	out  string
}

func (t *testTerm) Start() {
}

func (t *testTerm) Read(b []byte) (int, error) {
	a := <-t.keys
	copy(b, a)
	return len(a), nil
}

func (t *testTerm) Write(b []byte) (int, error) {
	// b = Deansi(b)
	t.out = t.out + string(b)
	return 1, nil
}

func (t *testTerm) Restore() error {
	return nil
}

func (t *testTerm) Close() error {
	close(t.keys)
	return nil
}

func (t *testTerm) RawMode() error {
	return nil
}
