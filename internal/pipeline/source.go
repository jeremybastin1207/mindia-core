package pipeline

import (
	"errors"
)

type ReadFunc = func(ctx PipelineCtx) (PipelineCtx, error)

type Source struct {
	read ReadFunc
}

func NewSource(read ReadFunc) Source {
	return Source{
		read,
	}
}

func (s *Source) Execute(ctx PipelineCtx) (PipelineCtx, error) {
	ctx, err := s.read(ctx)
	if ctx.Buffer == nil {
		return ctx, errors.New("source read function must return a buffer")
	}
	return ctx, err
}
