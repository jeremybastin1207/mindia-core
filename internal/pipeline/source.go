package pipeline

type GetterFunc = func(ctx PipelineCtx) (PipelineCtx, error)

type SourceConfig struct {
	Getter GetterFunc
}

type Source struct {
	SourceConfig
}

func NewSource(config SourceConfig) Source {
	return Source{
		SourceConfig: config,
	}
}

func (s *Source) Run(ctx PipelineCtx) (PipelineCtx, error) {
	return s.Getter(ctx)
}
