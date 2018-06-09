package ops

import (
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/util"
	"io"
	"strings"
)

type Op interface {
	Accept(str string) (bool, string)
	Complete() bool
	Completions(string) []string
	Hint(string) string
	Run(io.Reader, io.Writer, io.Writer)
	Str() string
}

type opIo struct {
	op  Op
	in  io.Reader
	out io.Writer
	err io.Writer
}

func newBuffer() *buffer {
	return &buffer{str: strings.Builder{}}
}

type buffer struct {
	str strings.Builder /// bytes.Buffer ?
}

func (b *buffer) Read(p []byte) (n int, err error) {
	s := []byte(b.str.String())
	copy(p, s)
	return len(s), nil
}

func (b *buffer) Write(s []byte) (int, error) {
	b.str.Write(s)
	return len(s), nil
}

var registry = map[string]func(*Ops) Op{}

func Init(args ...interface{}) {
}

func Register(name string, f func(*Ops) Op) {
	registry[name] = f
}

func New(conf *conf.Conf) *Ops {
	return &Ops{ops: []Op{}, Conf: conf}
}

type Ops struct {
	ops  []Op
	Conf *conf.Conf
}

func (o Ops) Run(in io.Reader, out io.Writer, err io.Writer) {
	if o.Complete() {
		for _, c := range chain(o.ops, in, out, err) {
			c.op.Run(c.in, c.out, c.err)
		}
	}
}

func (o *Ops) Accept(str string) bool {
	o.ops = []Op{}
	str = o.addOp(str, nil)
	for util.Strip(str) != "" {
		str = o.addOp(str, nil)
	}
	return len(o.ops) > 0 // this is always true once the first op is complete
}

func (o *Ops) addOp(str string, prev Op) string {
	for _, f := range registry {
		op := f(o)
		o.ops = append(o.ops, op)
		if ok, s := op.Accept(str); ok {
			return s
		} else {
			o.ops = o.ops[:len(o.ops)-1]
		}
	}

	op := api.New()
	if ok, str := op.Accept(str); ok {
		o.ops = append(o.ops, op)
		return str
	}

	return ""
}

func (o Ops) Complete() bool {
	if len(o.ops) == 0 {
		return false
	}
	ok := true
	for _, o := range o.ops {
		ok = ok && o.Complete()
	}
	return ok
}

func (o Ops) Completions(str string) []string {
	op := o.Last()
	if op != nil {
		return op.Completions(str)
	}
	op = api.New()
	return op.Completions(str)
}

func (o Ops) Hint(str string) string {
	op := o.Last()
	if op == nil {
		return ""
	}
	if op.Complete() && strings.HasSuffix(str, " ") {
		return "|"
	}
	strs := strings.Split(str, "|")
	str = strs[len(strs)-1]
	return op.Hint(strings.TrimPrefix(str, " "))
}

func (o Ops) Prev(op Op) Op {
	i := o.indexOf(op)
	if i > 0 {
		return o.ops[i-1]
	}
	return nil
}

func (o Ops) Len() int {
	return len(o.ops)
}

func (o Ops) First() Op {
	if len(o.ops) > 0 {
		return o.ops[0]
	}
	return nil
}

func (o Ops) Last() Op {
	if len(o.ops) > 0 {
		return o.ops[len(o.ops)-1]
	}
	return nil
}

func (o Ops) indexOf(op Op) int {
	for i, o := range o.ops {
		if o == op {
			return i
		}
	}
	return -1
}

func (o Ops) Str() string {
	strs := []string{}
	for _, o := range o.ops {
		strs = append(strs, o.Str())
	}
	return strings.Join(strs, "")
}

func chain(ops []Op, in io.Reader, out io.Writer, err io.Writer) []*opIo {
	c := []*opIo{}
	o := out
	if len(ops) > 1 {
		o = newBuffer()
	}
	for i, op := range ops {
		c = append(c, &opIo{op, in, o, err})
		in = c[i].out.(io.Reader)
		if i < len(ops)-2 {
			o = newBuffer()
		} else {
			o = out
		}
	}
	return c
}
