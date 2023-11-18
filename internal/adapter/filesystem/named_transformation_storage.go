package filesystem

import (
	"errors"
	"os"

	"github.com/jeremybastin1207/mindia-core/internal/transform"
	"gopkg.in/yaml.v2"
)

type NamedTransformationStorage struct {
	filename             string
	namedTransformations transform.NamedTransformationMap
}

func NewNamedTransformationStorage() *NamedTransformationStorage {
	return &NamedTransformationStorage{
		filename: "named_transformations.yml",
	}
}

func (s *NamedTransformationStorage) GetAll() (transform.NamedTransformationMap, error) {
	s.load()
	return s.namedTransformations, nil
}

func (s *NamedTransformationStorage) Get(name string) (*transform.NamedTransformation, error) {
	s.load()
	if val, ok := s.namedTransformations[name]; ok {
		return &val, nil
	}
	return nil, nil
}

func (s *NamedTransformationStorage) Save(namedTransformation transform.NamedTransformation) error {
	s.load()
	s.namedTransformations[namedTransformation.Name] = namedTransformation
	return s.save()
}

func (s *NamedTransformationStorage) Delete(name string) error {
	s.load()
	delete(s.namedTransformations, name)
	return s.save()
}

func (s *NamedTransformationStorage) DeleteAll() error {
	s.load()
	s.namedTransformations = map[string]transform.NamedTransformation{}
	return s.save()
}

func (s *NamedTransformationStorage) load() {
	if s.namedTransformations != nil {
		return
	}

	_, err := os.Stat(s.filename)
	if errors.Is(err, os.ErrNotExist) {
		err = s.save()
		if err != nil {
			panic(err)
		}
	}

	body, err := os.ReadFile(s.filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(body, &s.namedTransformations)
	if err != nil {
		panic(err)
	}
}

func (s *NamedTransformationStorage) save() error {
	yamlData, err := yaml.Marshal(s.namedTransformations)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, yamlData, 0644)
}
