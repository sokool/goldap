package main

import (
	"fmt"
	"github.com/sokool/console"
	"github.com/sokool/goldap"
	"github.com/sokool/goldap/filter"
)

func main() {
	ldap.Url = "michal.sokolowski:Haslo123@corp.test:389"
	adapter := ldap.New()

	f := filter.AND(
		filter.Contains("objectClass", "person"),
	)

	adapter.Search([]string{"*"}).When(f).Fetch().Each(func(i int, e *ldap.Element) {
		console.Log(fmt.Sprintf("%s [%s] ===> %s, %s", e.Value("displayName"), e.Value("location"), e.Value("description"), e.Value("comment")))
	})

}