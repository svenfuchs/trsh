package api

import (
	s "github.com/svenfuchs/trsh/spec"
	"github.com/svenfuchs/trsh/util"
	"strings"
)

func newResource(s *s.Resource) *resource {
	return &resource{
		spec: s,
		name: s.Name,
	}
}

type resource struct {
	spec   *s.Resource
	name   string
	action *action
}

func (r *resource) accept(str string) (bool, string) {
	str = util.StripPrefix(str, r.name)
	name := strings.Split(util.StripLeft(str), " ")[0]

	s := r.spec.FindAction(name)
	if s != nil {
		r.action = newAction(s)
	} else {
		r.action = nil
	}

	if r.action != nil {
		ok, str := r.action.accept(util.StripPrefix(str, r.action.name))
		return ok, str
	} else if strings.HasPrefix(r.name, str) {
		return true, ""
	} else if util.HasSuffix(str, " ") {
		return false, ""
	} else {
		return r.spec.MatchesAction(name), ""
	}
}

func (r *resource) hint(str string) string {
	if len(r.spec.ActionNames) == 0 {
		return ""
	}

	a := r.action
	if str == "" || a == nil {
		a = newAction(r.spec.Action(r.spec.ActionNames[0]))
	}

	s := strings.Split(str, " ")[0]
	h := a.hint(util.StripPrefix(str, s))
	return a.name + " " + h
}

func (r *resource) complete() bool {
	if r.action != nil {
		return r.action.complete()
	}
	return false
}

func (r resource) completions(str string) []string {
	str = util.StripLeft(str)

	if len(r.spec.ActionNames) == 0 {
		return []string{}
	}

	strs := r.spec.ActionNames
	if str != "" {
		strs = util.MatchStrs(strs, str)
	}
	if len(strs) == 1 {
		strs = []string{strs[0] + " "}
	}
	if r.action == nil || len(strs) > 0 {
		return strs
	}

	s := strings.TrimPrefix(str, r.action.name)
	if len(s) > 0 {
		return r.action.completions(util.Strip(s))
	}

	return []string{}
}

func (r resource) str() string {
	str := r.name
	if r.action != nil {
		str = str + " " + r.action.str()
	}
	return str
}

func (r resource) request() (string, string) {
	if r.action != nil {
		return r.action.request()
	}
	return "", ""
}
