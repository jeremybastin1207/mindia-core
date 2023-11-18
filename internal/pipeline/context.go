package pipeline

import "github.com/jeremybastin1207/mindia-core/internal/media"

type PipelineCtx struct {
	Path             media.Path
	Buffer           *Buffer
	ContentType      media.ContentType
	EmbeddedMetadata media.Metadata
	Tags             []media.Tag
}
