package types

type Anchor int

const (
	TopCenter Anchor = iota
	TopLeft
	TopRight
	BottomCenter
	BottomLeft
	BottomRight
	LeftCenter
	RightCenter
	Center
)
