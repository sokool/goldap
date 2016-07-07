package ldap

import (
	"gopkg.in/ldap.v2"
	"im.in/ldap/sanitizer"
)

type (
	Result struct {
		ldap.SearchResult
		filters map[string]func(*Element)
	}

	Element struct{ *ldap.Entry }
)

func (self *Element) Value(attribute string) string {
	return self.GetAttributeValue(attribute)
}

func (self *Element) SanitizeValue(method, attribute string) string {
	v, _ := sanitizer.Sanitize(method, []string{self.Value(attribute)}, []string{})

	return v[0]
}

func (self *Element) Each(fn func(name string, values []string, byteValues [][]byte)) {
	for _, attr := range self.Attributes {
		fn(attr.Name, attr.Values, attr.ByteValues)
	}
}

func (self *Result) RegisterFilter(name string, fn func(*Element)) {
	self.filters[name] = fn
}

func (self *Result) Count() int {
	return len(self.Entries)
}

func (self *Result) Each(fn func(int, *Element)) {
	for idx, e := range self.Entries {
		fn(idx, &Element{e})
	}
}
