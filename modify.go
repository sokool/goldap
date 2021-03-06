package ldap

import (
	"log"
	"gopkg.in/ldap.v2"
	"github.com/sokool/goldap/sanitizer"
)

type (
	Modify struct {
		sanitizer  *sanitizer.Sanitizer
		connection *LDAP
		dn         string
		toModify   *ldap.ModifyRequest
	}
)

func (m *Modify) In(dn string) *Modify {
	m.dn = dn
	return m
}

func (m *Modify) What(attr map[string][]string) *Modify {
	m.toModify = ldap.NewModifyRequest(m.dn)
	for name, values := range attr {
		if empty(values) {
			m.toModify.Delete(name, values)
			values = []string{" "}
		}

		for i, v := range values {
			o, _ := m.sanitizer.Sanitize(name, v)
			values[i] = o
		}

		m.toModify.Replace(name, values)
	}

	return m
}
func empty(in []string) bool {
	for _, v := range in {
		if len(v) != 0 {
			return false
		}
	}

	return true
}
func (m *Modify) Flush() error {
	m.connection.open()
	log.Printf("LDAP.Modify : %s", m.dn)
	err := m.connection.ldap.Modify(m.toModify)
	if err != nil {
		log.Printf("LDAP.Modify ERROR: %s", err.Error())
	}
	m.connection.close()
	return err
}
