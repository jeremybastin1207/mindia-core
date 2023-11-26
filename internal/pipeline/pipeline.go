package pipeline

import (
	"fmt"
)

type PipelineStep interface {
	Execute(ctx PipelineCtx) (PipelineCtx, error)
}

type Pipeline struct {
	steps []PipelineStep
}

func NewPipeline(source *Source, sinker *Sinker, steps []PipelineStep) Pipeline {
	pSteps := []PipelineStep{source}
	pSteps = append(pSteps, steps...)
	pSteps = append(pSteps, sinker)

	return Pipeline{
		steps: pSteps,
	}
}

func (p *Pipeline) Execute() (PipelineCtx, error) {
	var (
		pCtx PipelineCtx
		err  error
	)
	for i, t := range p.steps {
		pCtx, err = t.Execute(pCtx)
		if err != nil {
			return pCtx, fmt.Errorf("error executing step %d: %w", i, err)
		}
	}
	return pCtx, nil
}
