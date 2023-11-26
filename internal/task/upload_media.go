package task

import (
	"io"
	"time"

	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/parser"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

type UploadMediaTask struct {
	fileStorage               media.FileStorer
	cacheStorage              media.FileStorer
	mediaStorage              media.Storer
	mediaOptimization         transform.MediaOptimization
	namedTransformationParser parser.NamedTransformationParser
	transformationParser      parser.Parser
	transformationsBuilder    transform.Builder
}

func NewUploadMediaTask(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	mediaStorage media.Storer,
	namedTransformationStorage transform.Storer,
) UploadMediaTask {
	return UploadMediaTask{
		cacheStorage:              cacheStorage,
		fileStorage:               fileStorage,
		mediaStorage:              mediaStorage,
		mediaOptimization:         *transform.NewMediaOptimization(),
		namedTransformationParser: parser.NewNamedTransformationParser(namedTransformationStorage),
		transformationParser:      parser.NewParser(),
		transformationsBuilder:    transform.NewBuilder(fileStorage),
	}
}

func (t *UploadMediaTask) Upload(
	path string,
	body io.Reader,
	contentType string,
	contentLength int64,
	transformations []string,
) (*media.Media, error) {
	if !media.IsContentTypeSupported(contentType) {
		return nil, mindiaerr.New(mindiaerr.ErrCodeMimeTypeNotSupported)
	}

	source := pipeline.NewSource(
		func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
			ctx.Path = media.NewPath(path)
			ctx.Buffer = pipeline.NewBuffer(body)
			ctx.Buffer.ReadAll()
			ctx.ContentType = contentType
			return ctx, nil
		})

	sinker := pipeline.NewSinker(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
		err := t.fileStorage.Upload(media.UploadInput{
			Path:          ctx.Path,
			Body:          ctx.Buffer.Reader(),
			ContentType:   ctx.ContentType,
			ContentLength: ctx.Buffer.Len(),
		})
		return ctx, err
	})

	steps, err := t.mediaOptimization.GetSteps(contentType)
	if err != nil {
		return nil, err
	}
	p := pipeline.NewPipeline(&source, &sinker, steps)

	result, err := p.Execute()
	if err != nil {
		return nil, err
	}

	m := media.Media{
		Path:             result.Path,
		ContentType:      result.ContentType,
		ContentLength:    result.Buffer.Len(),
		EmbeddedMetadata: result.EmbeddedMetadata,
		DerivedMedias:    []media.DerivedMedia{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if len(transformations) > 0 {
		for _, transformation := range transformations {
			transformations, err := t.namedTransformationParser.Parse(transformation)
			if err != nil {
				return nil, err
			}
			trans, err := t.transformationParser.Parse(*transformations)
			if err != nil {
				return nil, err
			}
			steps, err := t.transformationsBuilder.Build(trans)
			if err != nil {
				return nil, err
			}

			source = pipeline.NewSource(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
				ctx.Path = result.Path
				ctx.Buffer = result.Buffer
				ctx.ContentType = result.ContentType
				ctx.EmbeddedMetadata = result.EmbeddedMetadata
				ctx.Buffer.ReadAll()
				return ctx, nil
			})

			cacheSinker := pipeline.NewSinker(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
				err := t.cacheStorage.Upload(media.UploadInput{
					Path:          ctx.Path.AppendSuffix(*transformations),
					Body:          ctx.Buffer.Reader(),
					ContentType:   ctx.ContentType,
					ContentLength: ctx.Buffer.Len(),
				})
				return ctx, err
			})

			p = pipeline.NewPipeline(&source, &cacheSinker, steps)
			_, err = p.Execute()
			if err != nil {
				return nil, err
			}
		}

		derivedMedias, err := t.cacheStorage.GetMultiple(m.Path)
		if err != nil {
			return nil, err
		}
		for _, a := range derivedMedias {
			m.DerivedMedias = append(m.DerivedMedias, media.DerivedMedia{
				Path:          a.Path,
				ContentType:   a.ContentType,
				ContentLength: a.ContentLength,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
		}
	}

	err = t.mediaStorage.Save(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
