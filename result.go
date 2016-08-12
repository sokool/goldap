package ldap

import (
	"gopkg.in/ldap.v2"
	"github.com/sokool/goldap/sanitizer"
	"encoding/json"
)

type (
	Result struct {
		ldap.SearchResult
		sanitizer *sanitizer.Sanitizer
	}

	Element struct {
		*ldap.Entry
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

func (r *Result) Count() int {
	return len(r.Entries)
}

func (r *Result) Each(fn func(int, *Element)) {
	for idx, e := range r.Entries {
		for _, a := range e.Attributes {
			r.sanitizer.Filter(a.Name, a)
		}

		fn(idx, &Element{e})

	}
}

func (r *Result) MarshalJSON() ([]byte, error) {
	out := []*Element{}
	r.Each(func(i int, e *Element) {
		out = append(out, e)
	})

	return json.Marshal(out)
}

