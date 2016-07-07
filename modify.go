package ldap

import (
	"gopkg.in/ldap.v2"
	"log"
	"github.com/sokool/console"
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
		console.Log(name, values)
		self.toModify.Replace(name, values)
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
