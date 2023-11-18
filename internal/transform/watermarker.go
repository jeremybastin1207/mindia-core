package transform

import (
	"bufio"
	"bytes"
	"image"

	"github.com/disintegration/imaging"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/pkg/types"
	"github.com/nickalie/go-webpbin"
)

type Watermark struct {
	Watermark media.Path
	Size      types.Size
	Anchor    types.Anchor
	Padding   string
}

type OverlaySinkerFunc = func() (media.Body, error)

type Watermarker struct {
	size          types.Size
	anchor        types.Anchor
	padding       int
	overlaySinker OverlaySinkerFunc
}

type WatermarkerConfig struct {
	Size          types.Size
	Anchor        types.Anchor
	Padding       int
	OverlaySinker OverlaySinkerFunc
}

func NewWatermarker(c WatermarkerConfig) Watermarker {
	return Watermarker{
		size:          c.Size,
		anchor:        c.Anchor,
		padding:       c.Padding,
		overlaySinker: c.OverlaySinker,
	}
}

func (w *Watermarker) Run(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
	dst, err := webpbin.Decode(ctx.Buffer.ReadAll())
	if err != nil {
		return ctx, err
	}

	overlayBody, err := w.overlaySinker()
	if err != nil {
		return ctx, err
	}
	overlay, err := webpbin.Decode(bytes.NewReader(*overlayBody))
	if err != nil {
		return ctx, err
	}

	if w.size.Width != 0 && w.size.Height != 0 {
		overlay = imaging.Fit(overlay, int(w.size.Width), int(w.size.Height), imaging.Lanczos)
	}

	pos := getWatermarkPosition(
		w.anchor,
		types.Size{
			Width:  int32(dst.Bounds().Dx()),
			Height: int32(dst.Bounds().Dy()),
		},
		types.Size{
			Width:  int32(overlay.Bounds().Dx()),
			Height: int32(overlay.Bounds().Dy()),
		},
		int32(w.padding),
	)

	dst = imaging.Overlay(dst, overlay, image.Pt(int(pos.X), int(pos.Y)), 1)

	err = webpbin.Encode(bufio.NewWriter(ctx.Buffer.ReadAll()), dst)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func getWatermarkPosition(anchor types.Anchor, dst types.Size, wk types.Size, padding int32) types.Position {
	var pos types.Position

	switch anchor {
	case types.TopLeft:
		pos.X = padding
		pos.Y = padding
	case types.TopCenter:
		pos.X = (dst.Width / 2) - (wk.Width / 2)
		pos.Y = padding
	case types.TopRight:
		pos.X = dst.Width - wk.Width - padding
		pos.Y = padding
	case types.BottomLeft:
		pos.X = padding
		pos.Y = dst.Height - wk.Height - padding
	case types.BottomCenter:
		pos.X = (dst.Width / 2) - (wk.Width / 2)
		pos.Y = dst.Height - wk.Height - padding
	case types.BottomRight:
		pos.X = dst.Width - wk.Width - padding
		pos.Y = dst.Height - wk.Height - padding
	case types.LeftCenter:
		pos.X = padding
		pos.Y = (dst.Height / 2) - (wk.Height / 2)
	case types.RightCenter:
		pos.X = dst.Width - wk.Width - padding
		pos.Y = (dst.Height / 2) - (wk.Height / 2)
	case types.Center:
		pos.X = (dst.Width / 2) - (wk.Width / 2)
		pos.Y = (dst.Height / 2) - (wk.Height / 2)
	default:
		pos.X = 0
		pos.Y = 0
	}

	return pos
}
