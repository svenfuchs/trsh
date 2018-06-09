package api

import (
	h "github.com/svenfuchs/trsh/http"
	s "github.com/svenfuchs/trsh/spec"
	"github.com/svenfuchs/trsh/util"
	"io"
	"strings"
)

var http *h.HTTP
var spec *s.Spec

func Init(args ...interface{}) {
	for _, arg := range args {
		switch a := arg.(type) {
		case *h.HTTP:
			http = a
		case *s.Spec:
			spec = a
		}
	}
	if http == nil {
		http = h.New()
	}
	if spec == nil {
		spec = s.Load(http)
	}
}

func New() *api {
	return &api{spec: spec}
}

type api struct {
	spec     *s.Spec
	resource *resource
}

func (a api) Run(in io.Reader, out io.Writer, err io.Writer) {
	req := a.request()
	if req != nil {
		util.Write(err, req.Method+" "+req.Path+" ...")
		res, _ := http.Run(req)
		util.Write(out, string(res.Body))
	}
}

func (a *api) Accept(str string) (bool, string) {
	name := strings.Split(str, " ")[0]

	if s := a.spec.FindResource(name); s != nil {
		a.resource = newResource(s)
	} else {
		a.resource = nil
	}

	if a.resource != nil {
		str := util.StripPrefix(str, a.resource.name)
		ok, str := a.resource.accept(str)
		return ok, str
	} else if util.HasSuffix(str, " ") {
		return false, ""
	}
	return a.spec.MatchesResource(name), ""
}

func (a *api) Complete() bool {
	if a.resource != nil {
		return a.resource.complete()
	}
	return false
}

func (a api) Completions(str string) []string {
	var strs []string
	if a.resource == nil {
		strs = a.spec.ResourceNames
	} else if !strings.HasPrefix(str, a.resource.name+" ") {
		strs = util.MatchStrs(a.spec.ResourceNames, str)
	} else {
		s := util.StripPrefix(str, a.resource.name)
		strs = a.resource.completions(s)
	}

	if len(strs) == 1 && strings.HasPrefix(strs[0], str) {
		strs = []string{strs[0] + " "}
	}
	return strs
}

func (a api) Hint(str string) string {
	r := a.resource
	if r == nil {
		r = newResource(a.spec.Resource(a.spec.ResourceNames[0]))
	}

	s := strings.Split(str, " ")[0]
	h := r.hint(util.StripPrefix(str, s))
	if !strings.HasPrefix(r.name+" "+h, str) {
		str = util.StripRight(str)
	}
	return strings.TrimPrefix(r.name+" "+h, str)
}

func (a api) request() *h.Request {
	if a.resource != nil && a.resource.complete() {
		return h.NewRequest(a.resource.request())
	}
	return nil
}

func (a api) Str() string {
	if a.resource != nil {
		return a.resource.str()
	}
	return ""
}
