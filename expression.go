package ldap


import (
	"fmt"
)

type (
	Expression string
)

func AND(f ... Expression) Expression {
	return through("&", f)
}

func OR(f ... Expression) Expression {
	return through("|", f)
}

func NOT(f ... Expression) Expression {
	return through("!", f)
}

func Equal(name, value string, def ...string) Expression {
	return parseDefault("(%s=%s)", name, value, def)
}

func Contains(name, value string, def ...string) Expression {
	return parseDefault("(%s=*%s*)", name, value, def)
}

func Gte(name, value string, def ...string) Expression {
	return parseDefault("(%s>=%s)", name, value, def)
}

func Lte(name, value string, def ...string) Expression {
	return parseDefault("(%s<=%s)", name, value, def)
}

func Between(name, start, end string) Expression {
	return Expression(string(Gte(name, start)) + string(Lte(name, end)))
}

func Present(v string) Expression {
	return parse()
}

func Ends(v string) Expression {
	return parse()
}

func Starts(v string) Expression {
	return parse()
}

func Approx(v string) Expression {
	return parse()
}

func Raw(v string) Expression {
	return parse()
}

func NewExpression(e ... Expression) Expression {
	return parse("(%s%s)", string(through("", e)), "")
}

func through(operator string, e []Expression) Expression {
	var expression Expression
	for _, x := range e {
		expression = expression + x
	}
	return parse("(%s%s)", operator, string(expression))
}

func parse(p ...string) Expression {
	return Expression(fmt.Sprintf(p[0], p[1], p[2]))
}

func parseDefault(pattern, name, value string, def []string) Expression {
	var v string
	if len(def) == 0 && value == "" {
		return Expression("")
	}

	if len(def) > 0 && value == "" {
		v = def[0]
	} else {
		v = value
	}

	return parse(pattern, name, v)
}
