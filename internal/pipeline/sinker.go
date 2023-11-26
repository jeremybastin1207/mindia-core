package pipeline

type SinkFunc = func(PipelineCtx) (PipelineCtx, error)

type Sinker struct {
	sink SinkFunc
}

func NewSinker(sink SinkFunc) Sinker {
	return Sinker{
		sink,
	}
}

func (s *Sinker) Execute(ctx PipelineCtx) (PipelineCtx, error) {
	return s.sink(ctx)
}
