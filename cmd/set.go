package cmd

import (
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/util"
	"io"
	"strings"
)

func newSet(conf *conf.Conf) ops.Op {
	return &set{conf: conf}
}

type set struct {
	name  string
	value string
	conf  *conf.Conf
}

func (s set) Run(io.Reader, io.Writer, io.Writer) {
	if s.name == "" {
		return
	} else if s.value == "" {
		s.conf.Del(s.name)
	} else {
		s.conf.Set(s.name, s.value)
	}
}

func (s set) Str() string {
	str := "set"
	if s.name != "" {
		str += " " + s.name + "=" + s.value
	}
	return str
}

func (s *set) Accept(str string) (bool, string) {
	s.name, s.value = "", ""
	strs := strings.Split(str, "=")
	name := util.Strip(strs[0])
	value := ""

	if len(strs) > 1 {
		value = util.Strip(strs[1])
	}

	if conf.IsOpt(name) {
		s.name = name
		s.value = value
		return true, ""
	}

	return conf.MatchesOpt(name), ""
}

func (s set) Complete() bool {
	return s.name != ""
}

func (s set) Completions(str string) []string {
	return []string{}
}

func (s set) Hint(str string) string {
	return ""
}
