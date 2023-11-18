package transform

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/nickalie/go-webpbin"
)

type WebpConverter struct {
}

func NewWebpConverter() WebpConverter {
	return WebpConverter{}
}

func (c *WebpConverter) Run(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
	var (
		img image.Image
		err error
	)

	if ctx.ContentType == string(media.ImageJpeg) || ctx.ContentType == string(media.ImageJpg) {
		img, err = jpeg.Decode(ctx.Buffer.ReadAll())
		if err != nil {
			return ctx, err
		}
	} else if ctx.ContentType == string(media.ImagePng) {
		img, err = png.Decode(ctx.Buffer.ReadAll())
		if err != nil {
			return ctx, err
		}
	} else if ctx.ContentType == string(media.ImageWebp) {
		return ctx, nil
	} else {
		return ctx, errors.New("unsupported content-type")
	}

	var buf bytes.Buffer
	err = webpbin.Encode(bufio.NewWriter(&buf), img)
	if err != nil {
		return ctx, err
	}

	ctx.Buffer.Body = buf.Bytes()
	ctx.Path = ctx.Path.SetExtension(".webp")
	ctx.ContentType = string(media.ImageWebp)

	return ctx, nil
}
