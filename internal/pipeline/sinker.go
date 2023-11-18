package pipeline

type SinkFunc = func(PipelineCtx) error

type Sinker struct {
	sinker SinkFunc
}

type SinkerConfig struct {
	Sinker SinkFunc
}

func NewSinker(c SinkerConfig) Sinker {
	return Sinker{
		sinker: c.Sinker,
	}
}

func (s *Sinker) Run(ctx PipelineCtx) (PipelineCtx, error) {
	err := s.sinker(ctx)
	return ctx, err
}
