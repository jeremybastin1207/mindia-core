package pipeline

type Step interface {
	Run(ctx PipelineCtx) (PipelineCtx, error)
}
