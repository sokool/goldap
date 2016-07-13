package main

import (
	"fmt"
	"github.com/sokool/goldap"
	"github.com/sokool/goldap/filter"
)

func main() {
	ldap.Url = "michal.sokolowski:dupa.romana.10@corp.smtsoftware.com:389"
	adapter := ldap.New()

	f := filter.AND(
		filter.Equal("objectClass", "person"),
	)

	adapter.Search([]string{"*"}).When(f).Fetch().Each(func(i int, e *ldap.Element) {
		fmt.Printf("%s %s\n", e.Value("displayName"), e.Value("l"))
		//console.Log(fmt.Sprintf("%s [%s] ===> %s, %s", e.Value("displayName"), e.Value("location"), e.Value("description"), e.Value("comment")))
	})

}