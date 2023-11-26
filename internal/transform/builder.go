package transform

import (
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
)

type Builder struct {
	scalerFactory      ScalerFactory
	watermarkerFactory WatermarkerFactory
}

func NewBuilder(dataStorage media.FileStorer) Builder {
	return Builder{
		scalerFactory:      ScalerFactory{},
		watermarkerFactory: NewWatermarkerFactory(dataStorage),
	}
}

func (b *Builder) Build(ts []Transformation) ([]pipeline.PipelineStep, error) {
	var transformers []pipeline.PipelineStep

	for _, t := range ts {
		var (
			t2  pipeline.PipelineStep
			err error
		)

		switch t.Name {
		case "c_scale":
			t2, err = b.scalerFactory.Build(t.Args)
		case "c_watermark":
			t2, err = b.watermarkerFactory.Build(t.Args)
		default:
			err = mindiaerr.New(mindiaerr.ErrCodeTransformationNotFound)
		}

		if err != nil {
			return nil, err
		}
		transformers = append(transformers, t2)
	}

	return transformers, nil
}
