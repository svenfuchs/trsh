package api

import (
	s "github.com/svenfuchs/trsh/spec"
)

func newParam(spec *s.Param, value string) *param {
	return &param{
		spec:  spec,
		name:  spec.Name,
		value: value,
		str:   spec.Name + "=" + string(value),
	}
}

type param struct {
	spec  *s.Param
	name  string
	value string
	str   string
}

func (p *param) complete() bool {
	return p.value != ""
}

// func (p param) str() string {
//   str := p.name + "="
//   if p.value != "" {
//     str = str + p.value
//   }
//   return str
// }
