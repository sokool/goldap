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
		if len(values) == 0 {
			m.toModify.Delete(name, values)
		} else {
			for i, v := range values {
				values[i], _ = m.sanitizer.Sanitize(name, v)
			}

			m.toModify.Replace(name, values)
		}
	}

	return m
}

func (m *Modify) Flush() error {
	m.connection.open()
	err := m.connection.ldap.Modify(m.toModify)
	if err != nil {
		log.Printf("LDAP.Modify ERROR: %s", err.Error())
	}
	m.connection.close()
	return err
}
