package transform

import (
	"image/color"
	"strconv"
	"strings"

	mindiacolor "github.com/jeremybastin1207/mindia-core/pkg/color"
	"github.com/jeremybastin1207/mindia-core/pkg/types"
)

type ScalerFactory struct {
}

func (f *ScalerFactory) Build(args map[string]string) (*Scaler, error) {
	w, _ := strconv.Atoi(args["w"])
	h, _ := strconv.Atoi(args["h"])

	opts := []ScalerOptions{
		WithSize(types.Size{Width: int32(w), Height: int32(h)}),
	}

	if args["a"] != "" {
		crop := ArgToCropStrategy(args["a"])
		opts = append(opts, WithCropStrategy(crop))
	}

	if args["s"] != "" {
		rgb, _ := mindiacolor.Hex2RGB(mindiacolor.Hex(strings.Replace(args["b"], "b_", "", 1)))
		opts = append(opts, WithPadColor(color.RGBA{rgb.Red, rgb.Green, rgb.Blue, 1}))
	}

	s := NewScaler(opts...)

	return &s, nil
}

func ArgToCropStrategy(arg string) CropStrategy {
	switch arg {
	case "pad":
		return PadResizeCrop
	case "forced":
		return ForcedCrop
	default:
		return ForcedCrop
	}
}
