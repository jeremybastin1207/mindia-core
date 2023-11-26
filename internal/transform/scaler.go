package transform

import (
	"bufio"
	"bytes"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/pkg/types"
	"github.com/nickalie/go-webpbin"
)

type CropStrategy string

var (
	PadResizeCrop CropStrategy = "pad_resize_crop"
	ForcedCrop    CropStrategy = "forced_crop"
)

type Scaler struct {
	size         types.Size
	cropStrategy CropStrategy
	padColor     color.RGBA
}

type ScalerConfig struct {
	Size         types.Size
	CropStrategy CropStrategy
	PadColor     color.RGBA
}

type ScalerOptions func(*scalerOptions)

type scalerOptions struct {
	size         types.Size
	cropStrategy CropStrategy
	padColor     color.RGBA
}

func newScalerOptions() *scalerOptions {
	return &scalerOptions{
		cropStrategy: CropStrategy(ForcedCrop),
		padColor:     color.RGBA{0, 0, 0, 1},
	}
}

func WithSize(size types.Size) ScalerOptions {
	return func(o *scalerOptions) {
		o.size = size
	}
}

func WithCropStrategy(cropStrategy CropStrategy) ScalerOptions {
	return func(o *scalerOptions) {
		o.cropStrategy = cropStrategy
	}
}

func WithPadColor(padColor color.RGBA) ScalerOptions {
	return func(o *scalerOptions) {
		o.padColor = padColor
	}
}

func NewScaler(opts ...ScalerOptions) Scaler {
	o := newScalerOptions()
	for _, optFunc := range opts {
		optFunc(o)
	}
	return Scaler{
		size:         o.size,
		cropStrategy: o.cropStrategy,
		padColor:     o.padColor,
	}
}

func (r *Scaler) Execute(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
	img, err := webpbin.Decode(ctx.Buffer.Reader())
	if err != nil {
		return ctx, err
	}

	img = imaging.Fit(img, int(r.size.Width), int(r.size.Height), imaging.Lanczos)
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()

	if r.cropStrategy == PadResizeCrop {
		if imgW != int(r.size.Width) || imgH != int(r.size.Height) {
			dst := imaging.New(int(r.size.Width), int(r.size.Height), color.RGBA{0, 0, 0, 0xff})
			pos := types.Position{X: 0, Y: 0}
			if img.Bounds().Dx() < int(r.size.Width) {
				pos.X = int32((int(r.size.Width) / 2) - (img.Bounds().Dx() / 2))
			} else {
				pos.Y = int32((int(r.size.Height) / 2) - (img.Bounds().Dy() / 2))
			}

			dst = imaging.Paste(dst, img, image.Pt(int(pos.X), int(pos.Y)))
			img = dst
		}
	}

	var buf bytes.Buffer
	err = webpbin.Encode(bufio.NewWriter(&buf), img)
	if err != nil {
		return ctx, err
	}

	ctx.Buffer = pipeline.NewBuffer(bytes.NewReader(buf.Bytes()))
	ctx.Path = ctx.Path.SetExtension(".webp")
	ctx.ContentType = string(media.ImageWebp)

	return ctx, nil
}
