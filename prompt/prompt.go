package prompt

import (
	"bytes"
	"github.com/svenfuchs/led-go"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/util"
	"io/ioutil"
	"os"
	"os/user"
)

func New(conf *conf.Conf, str string, t ...led.Iterm) *Prompt {
	e := led.NewEd(str, t...)
	p := &Prompt{ops: ops.New(conf), ed: e}
	e.Handle(led.Chars, p.insert)
	e.Handle(led.Backspace, p.backspace)
	e.Handle(led.Delete, p.delete)
	e.Handle(led.Tab, p.tab)
	e.Handle(led.ShiftTab, p.shiftTab)
	e.Handle(led.Left, p.left)
	e.Handle(led.Right, p.right)
	e.Handle(led.Up, p.prev)
	e.Handle(led.Down, p.next)
	e.Handle(led.Enter, p.enter)
	e.Handle(led.CtrlC, p.discard)
	return p
}

type Prompt struct {
	ops  *ops.Ops
	ed   *led.Ed
	strs [][]byte
}

func (p *Prompt) Run() {
	p.ed.Run()
}

func (p *Prompt) Str() string {
	return p.ed.Str()
}

func (p *Prompt) Suggested() string {
	return string(p.ed.Suggested)
}

func (p *Prompt) Pos() int {
	return p.ed.Pos
}

func (p *Prompt) insert(e *led.Ed, k led.Key) {
	str := e.Str()
	str = str[:e.Pos] + k.Str() + str[e.Pos:]

	if k.Str() == " " && util.HasSuffix(str, "  ", "= ") {
		e.Reject(k.Chars)
	} else if p.ops.Accept(str) {
		e.Insert(k.Chars)
	} else {
		e.Reject(k.Chars)
	}

	p.strs = nil
	p.hint(e)
}

func (p *Prompt) backspace(e *led.Ed, k led.Key) {
	p.ops.Accept(e.Str())
	e.Back()
	p.strs = nil
	p.hint(e)
}

func (p *Prompt) delete(e *led.Ed, k led.Key) {
}

func (p *Prompt) tab(e *led.Ed, k led.Key) {
	p.complete(e, led.Forw)
}

func (p *Prompt) shiftTab(e *led.Ed, k led.Key) {
	p.complete(e, led.Back)
}

func (p *Prompt) left(e *led.Ed, k led.Key) {
	p.backspace(e, k)
	p.hint(e)
}

func (p *Prompt) right(e *led.Ed, k led.Key) {
	str := e.Suggested
	if len(str) > 0 {
		e.Insert([]byte(str))
		p.ops.Accept(e.Str())
		p.hint(e)
	}
}

func (p *Prompt) enter(e *led.Ed, k led.Key) {
	historyAdd(e.Str())
	e.Pause()
	e.Newline()
	if p.ops.Accept(e.Str()) {
		p.ops.Run(os.Stdin, os.Stdout, os.Stderr)
	}
	e.Resume()
	e.Reset()
	p.reset()
}

func (p *Prompt) reset() {
	p.strs = nil
}

func (p *Prompt) discard(e *led.Ed, k led.Key) {
	e.Discard()
}

func (p *Prompt) complete(e *led.Ed, direction int) {
	str := e.Str()[:e.Pos]
	if p.strs == nil {
		strs := p.ops.Completions(str)
		p.strs = toByteArrays(strs)
	}
	if len(p.strs) > 0 {
		e.Complete(p.strs, direction)
		p.ops.Accept(e.Str())
		p.hint(e)
	}
	if util.HasSuffix(e.Str(), " ") {
		p.strs = nil
	}
}

func (p *Prompt) prev(e *led.Ed, k led.Key) {
	e.HistoryPrev(history())
	p.ops.Accept(e.Str())
}

func (p *Prompt) next(e *led.Ed, k led.Key) {
	e.HistoryNext(history())
	p.ops.Accept(e.Str())
}

func (p *Prompt) hint(e *led.Ed) {
	if e.Pos != len(e.Str()) {
		return
	}
	str := p.ops.Hint(e.Str())
	e.Suggest([]byte(str))
}

func historyFilename() string {
	u, _ := user.Current()
	return u.HomeDir + "/.travis-go/history"
}

func history() [][]byte {
	data, err := ioutil.ReadFile(historyFilename())
	if err == nil {
		return compact(bytes.Split(data, []byte("\n")))
	}
	return [][]byte{}
}

func historyAdd(line string) {
	if len(line) == 0 || contains(history(), line) {
		return
	}
	f, _ := os.OpenFile(historyFilename(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(line + "\n"))
	f.Close()
}

func toByteArrays(strs []string) [][]byte {
	b := [][]byte{}
	for _, s := range strs {
		b = append(b, []byte(s))
	}
	return b
}

func contains(s [][]byte, e string) bool {
	for _, a := range s {
		if string(a) == e {
			return true
		}
	}
	return false
}

func compact(lines [][]byte) [][]byte {
	r := [][]byte{}
	for _, l := range lines {
		if len(l) > 0 {
			r = append(r, l)
		}
	}
	return r
}
