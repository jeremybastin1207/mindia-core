package media

type Storer interface {
	Get(path Path) (*Media, error)
	GetMultiple(path Path, offset int, limit int, sortBy string, asc bool) ([]Media, error)
	Save(media *Media) error
	Delete(path Path) error
}
