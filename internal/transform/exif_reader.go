package transform

import (
	"fmt"
	"strings"

	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
)

type ExifReader struct {
}

func NewExifReader() ExifReader {
	exif.RegisterParsers(mknote.All...)
	return ExifReader{}
}

func (r *ExifReader) Execute(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
	var metadata = media.Metadata{}

	x, err := exif.Decode(ctx.Buffer.Reader())
	if err != nil {
		return ctx, nil
	}

	exifVersion, err := x.Get(exif.ExifVersion)
	if err == nil {
		metadata["exif_version"] = parseTag(exifVersion)
	}

	make, err := x.Get(exif.Make)
	if err == nil {
		metadata["make"] = parseTag(make)
	}

	model, err := x.Get(exif.Model)
	if err == nil {
		metadata["model"] = parseTag(model)
	}

	focal, err := x.Get(exif.FocalLength)
	if err == nil {
		numer, denom, err := focal.Rat2(0)
		if err == nil {
			metadata["focal_numerator"] = fmt.Sprintf("%v", numer)
			metadata["focal_denominator"] = fmt.Sprintf("%v", denom)
		}
	}

	tm, err := x.DateTime()
	if err == nil {
		metadata["taken"] = tm.String()
	}

	lat, long, err := x.LatLong()
	if err == nil {
		metadata["lat"] = fmt.Sprintf("%v", lat)
		metadata["long"] = fmt.Sprintf("%v", long)
	}

	colorSpace, err := x.Get(exif.ColorSpace)
	if err == nil {
		metadata["color_space"] = parseTag(colorSpace)
	}

	fNumber, err := x.Get(exif.FNumber)
	if err == nil {
		metadata["f_number"] = parseTag(fNumber)
	}

	ctx.EmbeddedMetadata = metadata

	return ctx, nil
}

func parseTag(t *tiff.Tag) string {
	return strings.ReplaceAll(t.String(), "\"", "")
}
