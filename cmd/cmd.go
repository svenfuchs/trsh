package cmd

import (
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/util"
	"io"
	"strings"
)

var cmds = map[string]func(*conf.Conf) ops.Op{
	"set": newSet,
}

func Init(_ ...interface{}) {
	ops.Register("cmd", New)
}

func Register(name string, factory func(*conf.Conf) ops.Op) {
	cmds[name] = factory
}

func New(ops *ops.Ops) ops.Op {
	return &dispatch{ops: ops}
}

type dispatch struct {
	cmd ops.Op
	ops *ops.Ops
}

func (d dispatch) Run(in io.Reader, out io.Writer, err io.Writer) {
	if d.cmd != nil {
		d.cmd.Run(in, out, err)
	}
}

func (d *dispatch) Accept(str string) (bool, string) {
	if !strings.HasPrefix(str, ":") {
		return false, ""
	}

	str = util.StripPrefix(str, ":")
	cmd := strings.Split(str, " ")[0]
	if cmd == "" {
		return true, ""
	}

	if c := d.findCmd(cmd); c != nil {
		d.cmd = c
	} else {
		d.cmd = nil
	}

	if d.cmd != nil {
		return d.cmd.Accept(util.StripPrefix(str, cmd))
	}
	return matchesCmd(cmd), ""
}

func (d dispatch) Complete() bool {
	if d.cmd != nil {
		return d.cmd.Complete()
	}
	return false
}

func (d dispatch) Completions(str string) []string {
	if d.cmd != nil {
		return d.cmd.Completions(str)
	}
	return []string{}
}

func (d dispatch) Hint(str string) string {
	if d.cmd != nil {
		return d.cmd.Hint(str)
	}
	return ""
}

func (d dispatch) Str() string {
	str := ":"
	if d.cmd != nil {
		str = str + d.cmd.Str()
	}
	return str
}

func (d dispatch) findCmd(cmd string) ops.Op {
	f := cmds[cmd]
	if f != nil {
		return f(d.ops.Conf)
	}
	return nil
}

func matchesCmd(str string) bool {
	for name, _ := range cmds {
		if strings.HasPrefix(name, str) {
			return true
		}
	}
	return false
}
