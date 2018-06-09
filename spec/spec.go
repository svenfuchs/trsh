package spec

import (
	"github.com/agnivade/levenshtein"
	"github.com/jinzhu/inflection"
	"github.com/jtacoma/uritemplates"
	"github.com/svenfuchs/trsh/http"
	"github.com/svenfuchs/trsh/util"
	"os/user"
	"sort"
	"strings"
)

// Load loads the api specification from a local file or from the api
func Load(h *http.HTTP, p ...string) *Spec {
	if len(p) == 0 {
		u, _ := user.Current()
		p = append(p, u.HomeDir+"/.trsh/manifest.json")
	}
	data := reader{http: h, path: p[0]}.run()
	return loadSpec(data)
}

// Spec represents the api specification
type Spec struct {
	Resources     map[string]*Resource
	ResourceNames []string
}

// Resource returns the resource with the given name
func (s Spec) Resource(name string) *Resource {
	return s.Resources[name]
}

// MatchesResource returns true if the given string matches a resource name
func (s Spec) MatchesResource(str string) bool {
	return util.MatchesStrs(s.ResourceNames, str)
}

// FindResource returns a resource matching the given string
func (s Spec) FindResource(str string) *Resource {
	// if str == "" {
	//   return nil
	// }
	m := util.MatchStrs(s.ResourceNames, str)
	if len(m) > 0 {
		return s.Resources[m[0]]
	}
	return nil
}

// FindSingularResource returns the singluar resource matching the given plural string
func (s Spec) FindSingularResource(str string) *Resource {
	inflection.AddIrregular("repository", "repositories")
	return s.FindResource(inflection.Singular(str))
}

// Suggest returns an array of suggestions for similar commands
func (s Spec) Suggest(str string) []string {
	m := []match{}
	for _, s := range s.Cmds() {
		m = append(m, match{s, levenshtein.ComputeDistance(str, s)})
	}
	sort.Sort(byDist{matches: m})
	strs := []string{}
	for _, match := range m[:10] {
		if match.dist == 1 {
			return []string{match.str}
		} else if match.dist < 4 {
			strs = append(strs, match.str)
		}
	}
	return strs
}

// Cmds returns an array of strings representing all known commands
func (s Spec) Cmds() []string {
	strs := []string{}
	for _, r := range s.Resources {
		strs = append(strs, r.Name)
		for _, a := range r.strs() {
			strs = append(strs, r.Name+" "+a)
		}
	}
	return strs
}

type match struct {
	str  string
	dist int
}

type byDist struct {
	matches []match
}

func (b byDist) Len() int {
	return len(b.matches)
}

func (b byDist) Swap(i, j int) {
	b.matches[i], b.matches[j] = b.matches[j], b.matches[i]
}

func (b byDist) Less(i, j int) bool {
	return b.matches[i].dist < b.matches[j].dist
}

// Resource represents a resource
type Resource struct {
	Name        string
	Actions     map[string]*Action
	ActionNames []string
	Attributes  []string
}

// Action returns the action with the given name
func (r *Resource) Action(name string) *Action {
	return r.Actions[name]
}

// HasAction returns true if the resource has an action with the given name
func (r *Resource) HasAction(name string) bool {
	return util.Includes(r.ActionNames, name)
}

// FindAction returns the action with the given name, if any
func (r *Resource) FindAction(str string) *Action {
	if str == "" {
		return nil
	}
	str = strings.Split(str, " ")[0]
	for _, n := range r.ActionNames {
		if n == str {
			return r.Action(n)
		}
	}
	return nil
}

// MatchesAction returns true if the given string matches an action name
func (r *Resource) MatchesAction(str string) bool {
	return util.MatchesStrs(r.ActionNames, str)
}

func (r *Resource) strs() []string {
	strs := []string{}
	for _, a := range r.Actions {
		strs = append(strs, a.strs()...)
	}
	return strs
}

// Action represents an action
type Action struct {
	Name   string
	Templs []*Templ
}

// FindTempl finds the template with params matching the given string
func (a *Action) FindTempl(str string) *Templ {
	if str == "" {
		return nil
	}
	str = strings.Split(str, "=")[0]
	for _, t := range a.Templs {
		if t.MatchesParam(str) {
			return t
		}
	}
	return nil
}

// MatchesParam returns true if the given string matches a param name
func (a *Action) MatchesParam(str string) bool {
	str = strings.Split(str, "=")[0]
	for _, t := range a.Templs {
		if t.MatchesParam(str) {
			return true
		}
	}
	return false
}

func (a *Action) strs() []string {
	strs := []string{}
	for _, t := range a.Templs {
		strs = append(strs, a.Name)
		for _, s := range t.strs() {
			strs = append(strs, a.Name+" "+s)
		}
	}
	return strs
}

// Templ represents an URI template
type Templ struct {
	Method       string
	Templ        *uritemplates.UriTemplate
	Params       []*Param
	ResourceName string
}

func newTempl(method string, templ *uritemplates.UriTemplate, resource string) *Templ {
	t := &Templ{
		Method:       method,
		Templ:        templ,
		Params:       []*Param{},
		ResourceName: resource,
	}
	for _, n := range util.Except(templ.Names(), "include") { // TODO
		t.Params = append(t.Params, newParam(resource, n))
	}
	return t
}

// MatchesParam returns true if the given string matches a param name
func (t *Templ) MatchesParam(str string) bool {
	return util.MatchesStrs(t.ParamNames(), str)
}

// Expand expands the template with the given params
func (t Templ) Expand(values map[string]string) string {
	v := make(map[string]interface{}) // no idea how to type cast to this
	for key, value := range values {
		v[key] = value
	}
	str, _ := t.Templ.Expand(v)
	return str
}

func (t Templ) strs() []string {
	p := t.ParamNames()
	if len(p) > 5 {
		p = p[:5]
	}
	s := []string{}
	for _, n := range p {
		s = append(s, n+"=")
	}
	return util.Combinations(s)
}

func (t Templ) ParamNames() []string {
	strs := []string{}
	for _, p := range t.Params {
		strs = append(strs, p.Name)
	}
	return strs
}

func (t Templ) ParamNamesEq() []string {
	strs := []string{}
	for _, p := range t.Params {
		strs = append(strs, p.Name+"=")
	}
	return strs
}

func newParam(resource string, name string) *Param {
	return &Param{
		FullName: name,
		Name:     strings.TrimPrefix(name, resource+"."),
	}
}

type Param struct {
	FullName string
	Name     string
}
