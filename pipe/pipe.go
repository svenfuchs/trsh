package pipe

import (
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/util"
	"io"
	"strings"
)

func Init(_ ...interface{}) {
	ops.Register("pipe", New)
}

func New(ops *ops.Ops) ops.Op {
	return &pipe{ops: ops}
}

type pipe struct {
	ops *ops.Ops
}

func (n pipe) Run(in io.Reader, out io.Writer, _ io.Writer) {
	util.Write(out, string(util.Read(in)))
}

func (n pipe) Str() string {
	return " | "
}

func (n pipe) Accept(str string) (bool, string) {
	str = util.StripLeft(str)
	if strings.HasPrefix(str, "|") {
		return true, util.StripPrefix(str, "|")
	}
	return false, ""
}

func (n pipe) Complete() bool {
	return n.ops.Len()%2 == 1
}

func (n pipe) Completions(str string) []string {
	return []string{}
}

func (n pipe) Hint(str string) string {
	return ""
}
