package parser

import (
	"strings"

	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

const transformationSeparator = "/"
const argSeparator = ","
const valueSeparator = "_"

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p *Parser) Parse(str string) ([]transform.Transformation, error) {
	var transformationsStr = []string{}

	for _, ts := range strings.Split(str, transformationSeparator) {
		transformationsStr = append(transformationsStr, ts)
	}

	var transformations = []transform.Transformation{}

	for _, el := range transformationsStr {
		name, args := parseTransformation(el)
		t := transform.NewTransformation(transform.Transformation{
			Name: name,
			Args: args,
		})
		transformations = append(transformations, t)
	}

	return transformations, nil
}

func parseTransformation(str string) (string, map[string]string) {
	var (
		name string
		args = make(map[string]string)
	)

	for i, arg := range strings.Split(str, argSeparator) {
		if i == 0 {
			name = arg
		} else {
			ap := strings.Split(arg, valueSeparator)
			args[ap[0]] = ap[1]
		}
	}
	return name, args
}
