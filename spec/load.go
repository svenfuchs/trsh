package spec

import (
	"github.com/jtacoma/uritemplates"
	"sort"
)

func loadSpec(data map[string]interface{}) *Spec {
	d := data["resources"].(map[string]interface{})
	r := loadResources(d)
	n := keys(d)
	return &Spec{Resources: r, ResourceNames: n}
}

func loadResources(data map[string]interface{}) map[string]*Resource {
	r := make(map[string]*Resource)
	for name, data := range data {
		r[name] = loadResource(name, data.(map[string]interface{}))
	}
	return r
}

func loadResource(name string, data map[string]interface{}) *Resource {
	d := data["actions"].(map[string]interface{})
	a := loadActions(name, d)
	n := sortActions(keys(d))
	s := toStrs(data["attributes"])
	return &Resource{Name: name, Actions: a, ActionNames: n, Attributes: s}
}

func loadActions(resource string, data map[string]interface{}) map[string]*Action {
	a := make(map[string]*Action)
	for name, data := range data {
		a[name] = loadAction(resource, name, data.([]interface{}))
	}
	return a
}

func sortActions(strs []string) []string {
	ix := indexOf(strs, "find")
	if ix == -1 {
		return strs
	}
	strs = append(strs[:ix], strs[ix+1:]...)
	return append([]string{"find"}, strs...)
}

func loadAction(resource string, name string, data []interface{}) *Action {
	return &Action{Name: name, Templs: loadTempls(resource, data)}
}

func loadTempls(resource string, data []interface{}) []*Templ {
	t := make([]*Templ, 0)
	for _, data := range data {
		t = append(t, loadTempl(resource, data.(map[string]interface{})))
	}
	return t
}

func loadTempl(resourceName string, data map[string]interface{}) *Templ {
	m := data["request_method"].(string)
	t := parseTempl(data["uri_template"].(string), resourceName)
	return newTempl(m, t, resourceName)
}

func parseTempl(str string, resourceName string) *uritemplates.UriTemplate {
	templ, _ := uritemplates.Parse(str)
	return templ
}

func keys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

func indexOf(strs []string, str string) int {
	for i, s := range strs {
		if s == str {
			return i
		}
	}
	return -1
}

func toStrs(data interface{}) []string {
	strs := []string{}
	for _, s := range data.([]interface{}) {
		strs = append(strs, s.(string))
	}
	return strs
}
