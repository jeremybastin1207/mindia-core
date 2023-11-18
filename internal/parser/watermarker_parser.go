package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
	"github.com/jeremybastin1207/mindia-core/pkg/types"
)

type WatermarkerFactory struct {
	dataStorage media.FileStorer
}

func NewWatermarkerFactory(dataStorage media.FileStorer) WatermarkerFactory {
	return WatermarkerFactory{
		dataStorage: dataStorage,
	}
}

func (f *WatermarkerFactory) Build(args map[string]string) (*transform.Watermarker, error) {
	w, _ := strconv.Atoi(args["w"])
	h, _ := strconv.Atoi(args["h"])
	p, _ := strconv.Atoi(args["p"])
	o := args["o"]
	o = strings.ReplaceAll(o, "@@", "/")

	a, err := stringToAnchor(args["a"])
	if err != nil {
		return &transform.Watermarker{}, errors.New("failed to parse anchor")
	}

	wm := transform.NewWatermarker(transform.WatermarkerConfig{
		OverlaySinker: func() (media.Body, error) {
			// body, _, err := f.dataStorage.Download(media.NewPath(o))
			// TODO return body, err
			return nil, nil
		},
		Size: types.Size{
			Width:  int32(w),
			Height: int32(h),
		},
		Anchor:  *a,
		Padding: p,
	})

	return &wm, nil
}

func stringToAnchor(s string) (*types.Anchor, error) {
	var a types.Anchor

	switch s {
	case "topcenter", "centertop":
		a = types.TopCenter
	case "topleft", "lefttop":
		a = types.TopLeft
	case "topright", "righttop":
		a = types.TopRight
	case "bottomcenter", "centerbottom":
		a = types.BottomCenter
	case "bottomleft", "leftbottom":
		a = types.BottomLeft
	case "bottomright", "rightbottom":
		a = types.BottomRight
	case "centerleft", "leftcenter":
		a = types.LeftCenter
	case "centerright", "rightcenter":
		a = types.RightCenter
	case "center":
		a = types.Center
	default:
		return nil, errors.New("invalid anchor String")
	}

	return &a, nil
}
