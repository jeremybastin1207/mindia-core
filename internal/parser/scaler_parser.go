package parser

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/jeremybastin1207/mindia-core/internal/transform"
	mindiacolor "github.com/jeremybastin1207/mindia-core/pkg/color"
	"github.com/jeremybastin1207/mindia-core/pkg/types"
)

type ScalerParser struct {
}

func (f *ScalerParser) Parse(args map[string]string) (*transform.Scaler, error) {
	w, _ := strconv.Atoi(args["w"])
	h, _ := strconv.Atoi(args["h"])

	opts := []transform.ScalerOptions{
		transform.WithSize(types.Size{Width: int32(w), Height: int32(h)}),
	}

	if args["a"] != "" {
		crop := ArgToCropStrategy(args["a"])
		opts = append(opts, transform.WithCropStrategy(crop))
	}

	if args["s"] != "" {
		rgb, _ := mindiacolor.Hex2RGB(mindiacolor.Hex(strings.Replace(args["b"], "b_", "", 1)))
		opts = append(opts, transform.WithPadColor(color.RGBA{rgb.Red, rgb.Green, rgb.Blue, 1}))
	}

	s := transform.NewScaler(opts...)

	return &s, nil
}

func ArgToCropStrategy(arg string) transform.CropStrategy {
	switch arg {
	case "pad":
		return transform.PadResizeCrop
	case "forced":
		return transform.ForcedCrop
	default:
		return transform.ForcedCrop
	}
}
