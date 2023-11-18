package transform

type Transformation struct {
	Name string
	Args map[string]string
}

func NewTransformation(c Transformation) Transformation {
	return Transformation{
		Name: c.Name,
		Args: c.Args,
	}
}

func (t *Transformation) IsSame(name string) bool {
	return t.Name == name
}

func (t *Transformation) ToString() string {
	return t.Name
}
