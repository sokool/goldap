package filter

import (
	"fmt"
)

var ()

type (
	Filter string
)

func AND(f ... Filter) Filter {
	return through("&", f)
}

func OR(f ... Filter) Filter {
	return through("|", f)
}

func NOT(f ... Filter) Filter {
	return through("!", f)
}

func Equal(name, value string, def ...string) Filter {
	return parseDefault("(%s=%s)", name, value, def)
}

func Contains(name, value string, def ...string) Filter {
	return parseDefault("(%s=*%s*)", name, value, def)
}

func Gte(name, value string, def ...string) Filter {
	return parseDefault("(%s>=%s)", name, value, def)
}

func Lte(name, value string, def ...string) Filter {
	return parseDefault("(%s<=%s)", name, value, def)
}

func Between(name, start, end string) Filter {
	return Filter(string(Gte(name, start)) + string(Lte(name, end)))
}

func Present(v string) Filter {
	return parse()
}

func Ends(v string) Filter {
	return parse()
}

func Starts(v string) Filter {
	return parse()
}

func Approx(v string) Filter {
	return parse()
}

func Raw(v string) Filter {
	return parse()
}

func New(e ... Filter) Filter {
	return parse("(%s%s)", string(through("", e)), "")
}

func through(operator string, e []Filter) Filter {
	var expression Filter
	for _, x := range e {
		expression = expression + x
	}
	return parse("(%s%s)", operator, string(expression))
}

func parse(p ...string) Filter {
	return Filter(fmt.Sprintf(p[0], p[1], p[2]))
}

func parseDefault(pattern, name, value string, def []string) Filter {
	var v string
	if len(def) == 0 && value == "" {
		return Filter("")
	}

	if len(def) > 0 && value == "" {
		v = def[0]
	} else {
		v = value
	}

	return parse(pattern, name, v)
}
