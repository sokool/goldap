package ldap

import "strings"

type tag struct {
	from    string
	action  string
	filters []string
}

func newTag(t string) tag {
	fs := strings.Split(t, ",")

	return tag{fs[0], fs[1], fs[2:]}

}
