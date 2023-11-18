package transform

import "time"

type NamedTransformationMap = map[string]NamedTransformation

const NamedTransformationPrefix = "t_"

type NamedTransformation struct {
	Name            string    `json:"name" yaml:"name"`
	Transformations string    `json:"transformations" yaml:"transformations"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
