package api

import (
	s "github.com/svenfuchs/trsh/spec"
	"github.com/svenfuchs/trsh/util"
	"sort"
	"strings"
)

func newAction(spec *s.Action) *action {
	return &action{
		spec: spec,
		name: spec.Name,
	}
}

type action struct {
	spec  *s.Action
	name  string
	templ *templ
}

func (a *action) accept(str string) (bool, string) {
	s := a.spec.FindTempl(str)
	if s != nil {
		a.templ = newTempl(s)
	}

	if a.templ != nil {
		ok, str := a.templ.accept(str)
		return ok, str
	} else if str == "" {
		return true, ""
	} else if strings.HasPrefix(a.name, str) {
		return true, strings.TrimPrefix(a.name, str)
	} else if util.HasSuffix(str, " ") {
		return false, ""
	} else {
		return a.spec.MatchesParam(str), ""
	}
}

func (a *action) hint(str string) string {
	if a.templ != nil {
		return a.templ.hint(str)
	}

	if len(a.paramNames()) > 0 {
		return a.paramNames()[0]
	}
	return ""
}

func (a *action) complete() bool {
	if a.templ != nil {
		return a.templ.complete()
	}
	return false
}

func (a action) completions(str string) []string {
	if a.templ == nil {
		return a.paramNames()
	}

	strs := []string{}
	for _, s := range a.spec.Templs {
		t := newTempl(s)
		if ok, _ := t.accept(str); ok {
			strs = append(strs, t.completions(str)...)
		}
	}

	sort.Strings(strs)
	return util.Unique(strs)
}

func (a *action) paramNames() []string {
	strs := []string{}
	m := map[string]bool{}
	for _, t := range a.spec.Templs {
		for _, n := range t.ParamNames() {
			if _, ok := m[n]; !ok {
				m[n] = true
				strs = append(strs, n+"=")
			}
		}
	}
	sort.Strings(strs)
	return util.Unique(strs)
}

func (a *action) request() (string, string) {
	if a.templ != nil {
		return a.templ.request()
	}
	return "", ""
}

func (a action) str() string {
	str := a.name
	if a.templ != nil {
		str = str + " " + a.templ.str()
	}
	return str
}
