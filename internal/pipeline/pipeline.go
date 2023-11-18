package pipeline

type Pipeline struct {
	steps []Step
}

type PipelineConfig struct {
	Source *Source
	Sinker *Sinker
	Steps  []Step
}

func NewPipeline(c PipelineConfig) Pipeline {
	steps := []Step{c.Source}
	steps = append(steps, c.Steps...)
	steps = append(steps, c.Sinker)
	return Pipeline{
		steps: steps,
	}
}

func (p *Pipeline) Execute() (PipelineCtx, error) {
	var (
		err error
		ctx PipelineCtx
	)
	for _, t := range p.steps {
		ctx, err = t.Run(ctx)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
