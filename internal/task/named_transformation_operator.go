package task

import (
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/transform"
	"github.com/jeremybastin1207/mindia-core/pkg/utils"
)

type NamedTransformationOperator struct {
	storage transform.Storer
}

func NewNamedTransformationOperator(storage transform.Storer) NamedTransformationOperator {
	return NamedTransformationOperator{
		storage,
	}
}

func (o *NamedTransformationOperator) GetAll() ([]transform.NamedTransformation, error) {
	t, err := o.storage.GetAll()
	if err != nil {
		return nil, err
	}
	return utils.ToArray(t), nil
}

func (o *NamedTransformationOperator) Create(
	name string,
	transformationsStr string,
) (*transform.NamedTransformation, error) {
	namedTransformation := transform.NamedTransformation{
		Name:            name,
		Transformations: transformationsStr,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := o.storage.Save(namedTransformation)
	if err != nil {
		return nil, err
	}
	return &namedTransformation, nil
}

func (o *NamedTransformationOperator) Update(
	name string,
	transformations string,
) (*transform.NamedTransformation, error) {
	t, err := o.storage.Get(name)
	if err != nil {
		return nil, err
	}
	t.Transformations = transformations
	err = o.storage.Save(*t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (o *NamedTransformationOperator) Delete(name string) error {
	return o.storage.Delete(name)
}

func (o *NamedTransformationOperator) DeleteAll() error {
	return o.storage.DeleteAll()
}
