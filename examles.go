package main

import (
	"im.in/ldap"
	"github.com/sokool/console"
	"im.in/ldap/filter"
	"fmt"
)

func main() {
	//"mail = %corp% AND (name = %a% OR name = "%b%") AND whenCreated >= 20160801 AND name != %siek%"
	ldap.Url = "michal.sokolowski:Haslo123@corp.test:389"
	adapter := ldap.New()

	f := filter.AND(
		filter.Contains("userPrincipalName", "corp"),
		//filter.OR(
		//filter.Contains("displayName", "michal"),
		//filter.Contains("displayName", "b"),
		//filter.Contains("displayName", "a"),
		//),
		//filter.Gte("whenCreated", "20160629063300.0Z"),
		//filter.NOT(
		//	filter.Contains("displayName", "siek"),
		//),
	)
	//f2 := filter.Expression("((&(displayName=*)(|(userPrincipalName=*a*)))")


	if user, exist := adapter.Search().When(filter.AND(filter.Contains("userPrincipalName", "sokol"))).FetchOne(); exist {
		m := map[string][]string{"displayName":{"Albert"}}
		adapter.Modify().In(user.DN).What(m).Flush()
	}

	adapter.Search().When(f).Fetch().Each(func(i int, e *ldap.Element) {
		console.Log(fmt.Sprintf("%s [%s] ===> %s, %s", e.Value("displayName"), e.Value("location"), e.Value("description"), e.Value("comment")))
	})

}