package media

import "io"

type UploadInput struct {
	Path          Path
	Body          io.Reader
	ContentType   ContentType
	ContentLength ContentLength
}

type DownloadResult struct {
	Path          Path
	Body          io.ReadCloser
	ContentType   ContentType
	ContentLength int64
}

type FileInfo struct {
	Path
	ContentType   ContentType   `json:"content_type,omitempty"`
	ContentLength ContentLength `json:"content_length,omitempty"`
}

type FileStorer interface {
	Upload(in UploadInput) error
	Download(p Path) (*DownloadResult, error)
	DownloadMultiple(p []Path) ([]*DownloadResult, error)
	Get(p Path) (*FileInfo, error)
	GetMultiple(p Path) ([]FileInfo, error)
	Move(src, dst Path) error
	Copy(src, dst Path) error
	Delete(p Path) error
	SpaceUsage() (int64, error)
}
