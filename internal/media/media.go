package media

import "time"

type Body = *[]byte
type ContentLength = int
type Metadata map[string]string

type Media struct {
	Path
	Body             Body           `json:"-"`
	ContentType      ContentType    `json:"content_type,omitempty"`
	ContentLength    ContentLength  `json:"content_length,omitempty"`
	EmbeddedMetadata Metadata       `json:"embedded_metadata,omitempty"`
	Tags             []Tag          `json:"tags,omitempty"`
	DerivedMedias    []DerivedMedia `json:"derived_medias,omitempty"`
	CreatedAt        time.Time      `json:"created_at,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty"`
}

type DerivedMedia struct {
	Path
	Body          Body          `json:"-"`
	ContentType   ContentType   `json:"content_type,omitempty"`
	ContentLength ContentLength `json:"content_length,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
}
