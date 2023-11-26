package task

import (
	"archive/zip"
	"bytes"
	"io"
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/analytics"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/parser"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

type DownloadMediaTask struct {
	fileStorage               media.FileStorer
	cacheStorage              media.FileStorer
	mediaStorage              media.Storer
	namedTransformationParser parser.NamedTransformationParser
	transformationParser      parser.Parser
	transformationsBuiler     transform.Builder
	analyticsRecorder         analytics.AnalyticsRecorder
}

func NewDownloadMediaTask(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	mediaStorage media.Storer,
	namedTransformationStorage transform.Storer,
	analyticsRecorder analytics.AnalyticsRecorder,
) DownloadMediaTask {
	return DownloadMediaTask{
		fileStorage:               fileStorage,
		cacheStorage:              cacheStorage,
		mediaStorage:              mediaStorage,
		namedTransformationParser: parser.NewNamedTransformationParser(namedTransformationStorage),
		transformationParser:      parser.NewParser(),
		transformationsBuiler:     transform.NewBuilder(fileStorage),
		analyticsRecorder:         analyticsRecorder,
	}
}

func (t *DownloadMediaTask) Download(
	transformations string,
	path media.Path,
) (io.ReadCloser, *media.ContentType, error) {
	var (
		downloadResult *media.DownloadResult
		err            error
	)

	defer func() {
		t.analyticsRecorder.RecordMediaRequest()
		if downloadResult != nil {
			t.analyticsRecorder.RecordBandwithUsage(int(downloadResult.ContentLength))
		}
	}()

	if transformations != "" {
		downloadResult, err = t.cacheStorage.Download(path.AppendSuffix(transformations))
		if err != nil {
			if err, ok := err.(*mindiaerr.Error); ok {
				if err.ErrCode != mindiaerr.ErrCodeMediaNotFound {
					return nil, nil, err
				}
			}
		} else {
			return downloadResult.Body, &downloadResult.ContentType, nil
		}
	} else {
		downloadResult, err = t.fileStorage.Download(path)
		if err != nil {
			return nil, nil, err
		}
		return downloadResult.Body, &downloadResult.ContentType, nil
	}

	parsedTransformations, err := t.namedTransformationParser.Parse(transformations)
	if err != nil {
		return nil, nil, err
	}
	trans, err := t.transformationParser.Parse(*parsedTransformations)
	if err != nil {
		return nil, nil, err
	}
	steps, err := t.transformationsBuiler.Build(trans)
	if err != nil {
		return nil, nil, err
	}

	source := pipeline.NewSource(
		func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
			downloadResult, err = t.fileStorage.Download(path)
			if err != nil {
				return ctx, err
			}
			ctx.Path = path
			ctx.ContentType = downloadResult.ContentType
			ctx.Buffer = pipeline.NewBuffer(downloadResult.Body)
			return ctx, nil
		})

	cacheSinker := pipeline.NewSinker(func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
		err := t.cacheStorage.Upload(media.UploadInput{
			Path:          ctx.Path.AppendSuffix(*parsedTransformations),
			Body:          ctx.Buffer.Reader(),
			ContentType:   ctx.ContentType,
			ContentLength: ctx.Buffer.Len(),
		})
		return ctx, err
	})

	p := pipeline.NewPipeline(&source, &cacheSinker, steps)

	result, err := p.Execute()
	if err != nil {
		return nil, nil, err
	}

	m, err := t.mediaStorage.Get(path)
	if err != nil {
		return nil, nil, err
	}

	derivedMedias, err := t.cacheStorage.GetMultiple(m.Path)
	if err != nil {
		return nil, nil, err
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

	err = t.mediaStorage.Save(m)
	if err != nil {
		return nil, nil, err
	}
	return io.NopCloser(result.Buffer.Reader()), &result.ContentType, nil
}

func (d *DownloadMediaTask) DownloadMultiple(paths []media.Path) (media.Body, error) {
	res, err := d.fileStorage.DownloadMultiple(paths)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	for _, r := range res {
		w, err := zipWriter.Create(r.Path.ToString()[1:])
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(w, r.Body)
		if err != nil {
			return nil, err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	bytes := buffer.Bytes()
	return &bytes, nil
}
