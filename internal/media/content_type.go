package media

import "golang.org/x/exp/slices"

type ContentType = string

const (
	ImageJpg  ContentType = "image/jpg"
	ImageJpeg ContentType = "image/jpeg"
	ImagePng  ContentType = "image/png"
	ImageWebp ContentType = "image/webp"
	VideoMp4  ContentType = "video/mp4"
	VideoMkv  ContentType = "video/x-matroska"
)

func IsContentTypeSupported(c ContentType) bool {
	cs := []ContentType{
		ImageJpeg,
		ImagePng,
		ImageWebp,
		VideoMp4,
		VideoMkv,
	}
	return slices.Contains(cs, c)
}
