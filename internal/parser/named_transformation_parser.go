package parser

import (
	"strings"

	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

const NamedTransformationPrefix = "t_"

type NamedTransformationParser struct {
	namedTransformationStorage transform.Storer
}

func NewNamedTransformationParser(namedTransformationStorage transform.Storer) NamedTransformationParser {
	if namedTransformationStorage == nil {
		panic("no storage provided")
	}
	return NamedTransformationParser{
		namedTransformationStorage,
	}
}

func (r *NamedTransformationParser) Parse(transformations string) (*string, error) {
	var transformationsStr = []string{}

	for _, ts := range strings.Split(transformations, transformationSeparator) {
		if strings.Contains(ts, NamedTransformationPrefix) {
			nt, err := r.namedTransformationStorage.Get(strings.Replace(ts, NamedTransformationPrefix, "", 1))
			if err != nil {
				return nil, err
			}
			transformationsStr = append(transformationsStr, strings.Split(nt.Transformations, transformationSeparator)...)
		} else {
			transformationsStr = append(transformationsStr, ts)
		}
	}

	transformations = strings.Join(transformationsStr[:], transformationSeparator)

	return &transformations, nil
}
