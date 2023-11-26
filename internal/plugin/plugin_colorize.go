package plugin

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
	"github.com/jeremybastin1207/mindia-core/pkg/utils"
	"github.com/replicate/replicate-go"
	"github.com/stretchr/objx"
)

const ColorizePluginName = "colorize"
const version = "376c74a2c9eb442a2ff9391b84dc5b949cd4e80b4dc0565115be0a19b7df0ae6"
const modelName = "Artistic"
const renderFactor = 35

type colorizeTaskDetails struct {
	Path         string `json:"path"`
	PredictionId string `json:"prediction_id"`
}

func NewColorizeTask(path media.Path) scheduler.Task {
	d := colorizeTaskDetails{
		Path:         path.ToString(),
		PredictionId: "",
	}
	return scheduler.NewTask(ColorizePluginName, d)
}

type ColorizePlugin struct {
	pluginManager   *PluginManager
	replicateClient *replicate.Client
}

func NewColorizePlugin(pluginManager *PluginManager) ColorizePlugin {
	replicateClient, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		utils.ExitErrorf("unable to create replicate client", err)
	}

	return ColorizePlugin{
		pluginManager,
		replicateClient,
	}
}

func (p *ColorizePlugin) Name() string {
	return ColorizePluginName
}

func (p *ColorizePlugin) Execute(task *scheduler.Task) (*scheduler.Task, error) {
	var d colorizeTaskDetails

	if reflect.TypeOf(task.Details) == reflect.TypeOf(map[string]interface{}{}) {
		d = colorizeTaskDetails{
			Path:         objx.New(task.Details).Get("path").Str(),
			PredictionId: objx.New(task.Details).Get("prediction_id").Str(),
		}
	} else {
		d = task.Details.(colorizeTaskDetails)
	}

	if d.Path == "" {
		return nil, errors.New("path is required")
	}

	path := media.NewPath(d.Path)

	if d.PredictionId == "" {
		prediction, err := p.createPrediction(path)
		if err != nil {
			return nil, err
		}
		details := colorizeTaskDetails{
			Path:         path.ToString(),
			PredictionId: prediction.ID,
		}
		task.Details = details
		task.Status = scheduler.Processing
		task.EnqueuedAt = time.Now()
	} else {
		url, status, err := p.fetchPrediction(path, d.PredictionId)
		if err != nil {
			return nil, err
		}
		if *status == replicate.Succeeded {
			err := p.savePicture(path, url)
			return nil, err
		}
		task.Status = toTaskStatus(*status)
	}

	return task, nil
}

func (p *ColorizePlugin) savePicture(path media.Path, url string) error {
	m, err := p.pluginManager.GetMediaStorage().Get(path)
	if err != nil {
		return err
	}

	source := pipeline.NewSource(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ctx, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ctx, err
		}
		defer resp.Body.Close()

		ctx.Path = m.Path
		ctx.Buffer = pipeline.NewBuffer(resp.Body)
		ctx.ContentType = m.ContentType
		ctx.EmbeddedMetadata = m.EmbeddedMetadata

		return ctx, nil
	})

	cacheSinker := pipeline.NewSinker(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
		err := p.pluginManager.GetCacheStorage().Upload(media.UploadInput{
			Path:          ctx.Path.AppendSuffix(ColorizePluginName),
			Body:          ctx.Buffer.Reader(),
			ContentType:   ctx.ContentType,
			ContentLength: ctx.Buffer.Len(),
		})
		return ctx, err
	})

	steps, err := p.pluginManager.GetMediaOptimization().GetSteps(media.ImageJpeg)
	if err != nil {
		return err
	}

	pp := pipeline.NewPipeline(&source, &cacheSinker, steps)
	_, err = pp.Execute()
	if err != nil {
		return err
	}

	derivedMedias, err := p.pluginManager.GetCacheStorage().GetMultiple(m.Path)
	if err != nil {
		return err
	}
	m.DerivedMedias = []media.DerivedMedia{}
	for _, a := range derivedMedias {
		m.DerivedMedias = append(m.DerivedMedias, media.DerivedMedia{
			Path:          a.Path,
			ContentType:   a.ContentType,
			ContentLength: a.ContentLength,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	err = p.pluginManager.GetMediaStorage().Save(m)
	if err != nil {
		return err
	}

	return nil
}

func (p *ColorizePlugin) createPrediction(path media.Path) (*replicate.Prediction, error) {
	media, err := p.pluginManager.GetFileStorage().Get(path)
	if err != nil {
		return nil, err
	}

	input := replicate.PredictionInput{
		"input_image":   "https://mindia-storage.ams3.digitaloceanspaces.com" + path.ToString(),
		"model_name":    modelName,
		"model_image":   media.Path.ToString(),
		"render_factor": renderFactor,
	}

	return p.replicateClient.CreatePrediction(context.Background(), version, input, nil, false)
}

func (p *ColorizePlugin) fetchPrediction(path media.Path, predictionId string) (string, *replicate.Status, error) {
	prediction, err := p.replicateClient.GetPrediction(context.Background(), predictionId)
	if err != nil {
		return "", nil, err
	}
	if prediction.Output != nil {
		predId := prediction.Output.(string)
		return predId, &prediction.Status, nil
	}
	return "", &prediction.Status, nil
}

func toTaskStatus(s replicate.Status) scheduler.TaskStatus {
	switch s {
	case replicate.Starting, replicate.Processing:
		return scheduler.Processing
	case replicate.Failed, replicate.Succeeded:
		return scheduler.Finished
	default:
		return scheduler.Canceled
	}
}

func (p *ColorizePlugin) getTaskName(path media.Path) string {
	return p.Name() + ":" + path.ToString()
}
