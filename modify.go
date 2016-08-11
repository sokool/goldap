package ldap

import (
	"log"

	"gopkg.in/ldap.v2"
)

type (
	Modify struct {
		connection *LDAP
		dn         string
		toModify   *ldap.ModifyRequest
	}
)

func (self *Modify) In(dn string) *Modify {
	self.dn = dn
	return self
}

func (self *Modify) What(attrs map[string][]string) *Modify {
	self.toModify = ldap.NewModifyRequest(self.dn)

	for name, values := range attrs {
		if len(values) == 0 {
			self.toModify.Delete(name, values)
		} else {
			self.toModify.Replace(name, values)
		}
	}

	return self
}

func (self *Modify) Flush() error {
	self.connection.open()
	err := self.connection.ldap.Modify(self.toModify)
	if err != nil {
		log.Printf("LDAP.Modify ERROR: %s", err.Error())
	}
	self.connection.close()
	return err
}
