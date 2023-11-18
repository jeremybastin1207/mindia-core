package plugin

import (
	"context"
	"net/http"
	"time"

	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
	"github.com/replicate/replicate-go"
	"github.com/stretchr/objx"
)

const version = "376c74a2c9eb442a2ff9391b84dc5b949cd4e80b4dc0565115be0a19b7df0ae6"
const task_name = "colorize"

type ColorizeTaskDetails struct {
	Path         string `json:"path"`
	PredictionId string `json:"prediction_id"`
}

type ColorizePlugin struct {
	pluginManager   *PluginManager
	replicateClient *replicate.Client
}

func NewColorizePlugin(pluginManager *PluginManager) ColorizePlugin {
	replicateClient, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		mindiaerr.ExitErrorf("unable to create replicate client", err)
	}

	return ColorizePlugin{
		pluginManager,
		replicateClient,
	}
}

func (p *ColorizePlugin) Name() string {
	return "colorize"
}

func (p *ColorizePlugin) Execute(path media.Path) error {
	media, err := p.pluginManager.GetFileStorage().Get(path)
	if err != nil {
		return err
	}

	input := replicate.PredictionInput{
		"input_image":   "https://mindia-storage.ams3.digitaloceanspaces.com" + path.ToString(),
		"model_name":    "Artistic",
		"model_image":   media.Path.ToString(),
		"render_factor": 35,
	}

	prediction, err := p.replicateClient.CreatePrediction(context.Background(), version, input, nil, false)
	if err != nil {
		return err
	}

	details := ColorizeTaskDetails{
		Path:         media.Path.ToString(),
		PredictionId: prediction.ID,
	}

	task := scheduler.NewTask(task_name, details)
	task.Status = scheduler.Processing
	task.EnqueuedAt = time.Now()
	err = p.pluginManager.GetTaskStorage().Save(task)
	if err != nil {
		return err
	}
	return nil
}

func (p *ColorizePlugin) Hook(task scheduler.Task) error {
	d := ColorizeTaskDetails{
		Path:         objx.New(task.Details).Get("path").Str(),
		PredictionId: objx.New(task.Details).Get("prediction_id").Str(),
	}

	prediction, err := p.replicateClient.GetPrediction(context.Background(), d.PredictionId)
	if err != nil {
		return err
	}

	switch prediction.Status {
	case replicate.Succeeded, replicate.Failed, replicate.Canceled:
		task.Status = scheduler.Finished
	default:
		return nil
	}

	m, err := p.pluginManager.GetMediaStorage().Get(media.NewPath(d.Path))
	if err != nil {
		task.Status = scheduler.Canceled
	}

	source := pipeline.NewSource(pipeline.SourceConfig{
		Getter: func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
			ctx.Path = media.NewPath(d.Path)

			req, _ := http.NewRequest("GET", prediction.Output.(string), nil)
			resp, _ := http.DefaultClient.Do(req)

			ctx.Buffer = &pipeline.Buffer{
				Reader: resp.Body,
			}
			ctx.ContentType = m.ContentType
			return ctx, nil
		},
	})

	cacheSinker := pipeline.NewSinker(pipeline.SinkerConfig{
		Sinker: func(ctx pipeline.PipelineCtx) error {
			ctx.Path = ctx.Path.AppendSuffix("colorize")

			return p.pluginManager.GetCacheStorage().Upload(media.UploadInput{
				Path:          ctx.Path,
				Body:          ctx.Buffer.MergeReader(),
				ContentType:   ctx.ContentType,
				ContentLength: ctx.Buffer.Len(),
			})
		},
	})

	pp := pipeline.NewPipeline(pipeline.PipelineConfig{
		Source: &source,
		Sinker: &cacheSinker,
	})

	result, err := pp.Execute()
	if err != nil {
		return err
	}

	m.DerivedMedias = append(m.DerivedMedias, media.DerivedMedia{
		Path:          result.Path,
		ContentType:   result.ContentType,
		ContentLength: result.Buffer.Len(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})

	err = p.pluginManager.GetMediaStorage().Save(m)
	if err != nil {
		return err
	}

	return p.pluginManager.taskStorage.Save(task)
}
