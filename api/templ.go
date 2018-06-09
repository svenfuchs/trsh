package api

import (
	s "github.com/svenfuchs/trsh/spec"
	"github.com/svenfuchs/trsh/util"
	"regexp"
	"sort"
	"strings"
)

func newTempl(spec *s.Templ) *templ {
	return &templ{
		spec:   spec,
		params: []*param{},
	}
}

type templ struct {
	spec   *s.Templ
	params []*param
}

func (t *templ) accept(str string) (bool, string) {
	t.params = []*param{}
	s := util.Strip(strings.Split(str, "|")[0])
	strs := strings.Split(str, " ")

	for _, s := range strs {
		parts := strings.Split(s, "=")
		name := parts[0]

		if name == "" {
			continue
		}

		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}

		for _, s := range t.spec.Params {
			if strings.HasPrefix(s.Name, name) {
				t.params = append(t.params, newParam(s, value))
				break
			}
		}
	}

	var ok = true
	for _, s := range strings.Split(s, " ") {
		if !t.spec.MatchesParam(strings.Split(s, "=")[0]) {
			ok = false
			break
		}
	}
	return ok, strings.TrimPrefix(str, s)
}

func (t *templ) complete() bool {
	for _, p := range t.params {
		if p.value == "" {
			return false
		}
	}
	return true
}

func (t templ) completions(str string) []string {
	str = util.LastPart(str, " ")
	strs := util.MatchStrs(t.spec.ParamNamesEq(), str)
	strs = util.Except(strs, str)
	sort.Strings(strs)
	return util.Unique(strs)
}

var matchParam = regexp.MustCompile(`[\w\.]+=[\w]+\s?`)

func (t *templ) hint(str string) string {
	if str == "" && len(t.spec.ParamNamesEq()) == 0 {
		return t.str()
	} else if str == "" {
		return t.spec.ParamNamesEq()[0]
	}

	s := matchParam.ReplaceAllString(str, "")
	if s == "" && !strings.HasSuffix(str, " ") {
		return t.str()
	}

	strs := []string{}
	for _, n := range t.spec.ParamNames() {
		if p := t.param(n); p == nil || !p.complete() {
			strs = append(strs, n+"=")
		}
	}
	if s != "" {
		strs = util.MatchStrs(strs, s)
	}
	str = t.str()

	if len(strs) > 0 && !strings.Contains(str, strs[0]) {
		return strings.TrimPrefix(str+" "+strs[0], " ")
	}

	return str
}

func (t *templ) param(name string) *param {
	for _, p := range t.params {
		if p.spec.Name == name {
			return p
		}
	}
	return nil
}

func (t templ) str() string {
	strs := []string{}
	for _, p := range t.params {
		strs = append(strs, p.str)
	}
	return strings.Join(strs, " ")
}

func (t *templ) request() (string, string) {
	v := make(map[string]string)
	for _, p := range t.params {
		v[p.spec.FullName] = p.value
	}
	return t.spec.Method, t.spec.Expand(v)
}
