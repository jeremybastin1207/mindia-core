package transform

type Storer interface {
	GetAll() (NamedTransformationMap, error)
	Get(name string) (*NamedTransformation, error)
	Save(t NamedTransformation) error
	Delete(name string) error
	DeleteAll() error
}
