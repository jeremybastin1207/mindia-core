package transform

import (
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
)

func NewGoogleTagger() GoogleTagger {
	return GoogleTagger{}
}

type GoogleTagger struct {
}

func (w *GoogleTagger) Run(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
	var (
		ctx2 = context.Background()
	)

	client, err := vision.NewImageAnnotatorClient(ctx2)
	if err != nil {
		return ctx, err
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(ctx.Buffer.Reader())
	if err != nil {
		return ctx, err
	}

	labels, err := client.DetectLabels(ctx2, image, nil, 10)
	if err != nil {
		return ctx, err
	}

	var tags []media.Tag
	for _, label := range labels {
		tags = append(tags, media.Tag{
			Value:           label.Description,
			ConfidenceScore: label.Score,
			Provider:        "google",
		})
	}

	ctx.Tags = tags

	return ctx, nil
}
