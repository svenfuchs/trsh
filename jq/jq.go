package jq

import (
	Jq "github.com/savaki/jq"
	"github.com/svenfuchs/trsh/http"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/spec"
	"github.com/svenfuchs/trsh/util"
	"io"
	"regexp"
	"strings"
)

var s *spec.Spec

func Init(args ...interface{}) {
	for _, arg := range args {
		switch a := arg.(type) {
		case *spec.Spec:
			s = a
		}
	}
	if s == nil {
		s = spec.Load(http.New())
	}
	ops.Register("jq", New)
}

func New(ops *ops.Ops) ops.Op {
	return &jq{spec: s, ops: ops, query: ""}
}

type jq struct {
	ops   *ops.Ops
	spec  *spec.Spec
	query string
}

func (j jq) Run(in io.Reader, out io.Writer, err io.Writer) {
	q, _ := Jq.Parse(j.query)
	b, _ := q.Apply(util.Read(in))
	util.Write(out, string(b))
}

func (j *jq) Accept(str string) (bool, string) {
	j.query = ""

	if !strings.HasPrefix(str, ".") || !j.acceptAttrs(str) {
		return false, ""
	}

	if _, err := Jq.Parse(str); err == nil {
		j.query = str
		return true, ""
	}

	return false, ""
}

func (j *jq) acceptAttrs(str string) bool {
	s := j.spec.Resource(j.resource())
	if s == nil {
		return false
	}

	strs := parseJq(str)
	if len(strs) > 2 {
		j.query = ""
		return false
	}

	if len(strs) > 0 {
		for _, str := range strs {
			if util.FindMatch(s.Attributes, str) == "" {
				j.query = ""
				return false
			}
			s = j.spec.FindSingularResource(j.resource())
		}
	}

	return true
}

func (j jq) Complete() bool {
	return true
}

func (j *jq) Hint(str string) string {
	s := j.spec.Resource(j.resource())
	if s == nil {
		return ""
	}

	if str == "" {
		return "."
	}
	str = strings.TrimPrefix(str, ".")
	strs := s.Attributes
	name := util.FirstOr(strs, "")
	if str != "" {
		name = util.FindMatch(strs, str)
	}
	return strings.TrimPrefix(name, str)
}

func (j *jq) Completions(str string) []string {
	s := j.spec.Resource(j.resource())
	if s == nil {
		return []string{}
	}

	str = strings.TrimPrefix(str, ".")
	strs := s.Attributes
	if str != "" {
		strs = util.MatchStrs(strs, str)
	}
	return strs
}

func (j jq) Str() string {
	return j.query
}

func (j *jq) resource() string {
	op := j.ops.Prev(j)
	if op == nil {
		return ""
	}
	op = j.ops.Prev(op)
	if op == nil {
		return ""
	}
	return strings.Split(op.Str(), " ")[0]
}

var jqRegex = []*regexp.Regexp{
	regexp.MustCompile(`\.([\w]+)`),
	regexp.MustCompile(`\[\s*"(.*)"\s*\]`),
}

func parseJq(str string) []string {
	keys := []string{}
	for _, r := range jqRegex {
		m := r.FindAllStringSubmatch(str, -1)
		if len(m) > 0 {
			for _, a := range m {
				keys = append(keys, a[1])
			}
			continue
		}
	}
	return keys
}
