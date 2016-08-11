package ldap

import (
	"gopkg.in/ldap.v2"
	"github.com/sokool/goldap/sanitizer"
	"encoding/json"
)

type (
	Result struct {
		ldap.SearchResult
		Filters map[string][]string
	}

	Element struct {
		*ldap.Entry
		filters map[string][]string
	}
)
//
//func (e *Element) Value(attribute string) string {
//	return e.GetAttributeValue(attribute)
//}
//
//func (e *Element) SanitizeValue(method, attribute string) string {
//	v, _ := sanitizer.Sanitize(method, []string{e.Value(attribute)}, []string{})
//
//	return v[0]
//}

func (e *Element) Each(fn func(*ldap.EntryAttribute)) {
	for _, a := range e.Attributes {
		if f, ok := e.filters[a.Name]; ok {
			for _, m := range f {
				sanitizer.Run(m, a)
			}
		}
		fn(a)
	}
}

func (e *Element) MarshalJSON() ([]byte, error) {
	out := map[string][]string{}
	e.Each(func(a *ldap.EntryAttribute) {
		for _, v := range a.Values {
			out[a.Name] = append(out[a.Name], v)
		}
	})

	return json.Marshal(out)
}

func (self *Result) RegisterFilter(filter string, fields []string) *Result {
	for _, n := range fields {
		self.Filters[n] = append(self.Filters[n], filter)
	}

	return self
}

func (self *Result) Count() int {
	return len(self.Entries)
}

func (self *Result) Each(fn func(int, *Element)) {
	for idx, e := range self.Entries {
		fn(idx, &Element{e, self.Filters})
	}
}

func (e *Result) MarshalJSON() ([]byte, error) {
	out := []*Element{}
	e.Each(func(i int, e *Element) {
		out = append(out, e)
	})

	return json.Marshal(out)
}

