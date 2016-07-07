package ldap

import (
	"gopkg.in/ldap.v2"
	"log"
	"strings"
	"errors"
	"fmt"
)

type (
	LDAP struct {
		url  string
		ldap *ldap.Conn
	}
)

var (
	Url, username, password, domain, port string
	dns []dn
)

func (self *LDAP) close() {
	log.Printf("LDAP.Close")
	self.ldap.Close()
}

func (self *LDAP) open() {
	var err error

	log.Printf("LDAP.Dial: %s:%s", domain, port)
	self.ldap, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", domain, port))
	if err != nil {
		log.Fatal(err)
	}

	user := fmt.Sprintf("%s@%s", username, domain)
	log.Printf("LDAP.Bind: %s:%s", user, password)

	err = self.ldap.Bind(user, password)
	if err != nil {
		log.Fatal(err)
	}
}

func (self *LDAP) Authenticate(usr, pass string) bool {
	var err error

	if usr == "" || pass == "" {
		log.Printf("LDAP.Authentication : empty user or pass variables")
		return false
	}

	self.ldap, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", domain, port))
	if err != nil {
		log.Fatal(err)
	}

	//check if credentials are OK
	u := fmt.Sprintf("%s@%s", usr, domain)
	log.Printf("LDAP.Authentication of: %s", u)
	err = self.ldap.Bind(u, pass)
	if err != nil {
		log.Println("LDAP.Authentication FAILED: %s", err)
		return false
	}

	return true
}

func (self *LDAP) Search(attributes []string) *Search {
	return &Search{
		connection: self,
		domain: dns,
		attributes: attributes,
	}
}

func (self *LDAP) Modify() *Modify {
	return &Modify{
		connection:self,
	}
}

func (self *LDAP) Add() {

	//CN=Schema,CN=Configuration,DC=corp,DC=test
	//changetype: modify
	//add: attributeTypes
	//attributeTypes: ( 1.3.6.1.4.1.32473.1.1.590
	//NAME ( 'blog' 'blogURL' )
	//DESC 'URL to a personal weblog'
	//SYNTAX 1.3.6.1.4.1.1466.115.121.1.15
	//SINGLE-VALUE
	//X-ORIGIN 'Oracle Unified Directory Server'
	//USAGE userApplications )

}

func (self *LDAP) Delete() {

}

func New() *LDAP {
	var err error
	username, password, domain, port, dns, err = parseURL(Url)
	if err != nil {
		log.Println(err)
	}

	return &LDAP{
		url: Url,
	}
}

func parseURL(url string) (string, string, string, string, []dn, error) {
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		return "", "", "", "", nil, errors.New("Invalid username:password@host:port format")
	}

	uParts := strings.Split(parts[0], ":")
	if len(uParts) != 2 {
		return "", "", "", "", nil, errors.New("Invalid username:password format")
	}

	dParts := strings.Split(parts[1], ":")
	if len(dParts) != 2 {
		return "", "", "", "", nil, errors.New("Invalid host:port format")
	}

	return uParts[0], uParts[1], dParts[0], dParts[1], toDN(dParts[0], ".", "dc"), nil
}
