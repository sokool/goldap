package ldap

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"strings"
	"log"
	f "github.com/sokool/goldap/filter"
)

type (
	Search struct {
		connection *LDAP
		domain     []dn
		dns        []dn
		filters    f.Filter
		attributes []string
	}

	dn struct {
		name, value string
	}
)

func (self *dn) String() string {
	return self.name + "=" + self.value
}

func (self *Search) In(namespace, value string) *Search {
	self.dns = append(self.dns, dn{namespace, value})
	return self
}

func (self *Search) When(filter f.Filter) *Search {
	self.filters = filter
	return self
}

func (self *Search) proceed() *Result {
	var bases []string
	for _, b := range self.dns {
		bases = append(bases, fmt.Sprintf("%s=%s", b.name, b.value))
	}

	base := strings.Join(bases, ",")

	//log.Printf("Select.BaseDN: %s\n", base)
	//log.Printf("Select.Filters: %s\n", string(self.filters))
	//log.Printf("Select.Attributes: %s\n", self.attributes)

	request := ldap.NewSearchRequest(
		base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 5000, 0, false,
		string(self.filters),
		self.attributes,
		nil,
	)

	response, err := self.connection.ldap.Search(request)
	log.Printf("LDAP.Found: %d items, query %s\n", len(response.Entries), string(self.filters))
	if err != nil {
		log.Printf("Select.Search ERROR: %s", err)
	}

	var result Result
	result = Result{SearchResult: *response}
	return &result
}

func (self *Search) Fetch() *Result {
	self.dns = append(self.dns, self.domain...)

	self.connection.open()

	series := self.proceed()

	self.connection.close()

	return series
}

func (self *Search) FetchOne() (*Element, bool) {
	result := self.Fetch()
	if len(result.Entries) != 1 {
		return nil, false
	}

	return &Element{result.Entries[0]}, true
}

func toDN(query, sep, namespace string) []dn {
	var dns []dn
	for _, v := range strings.Split(query, sep) {
		dns = append(dns, dn{namespace, v})
	}
	return dns
}