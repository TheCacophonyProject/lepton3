package lepton3

import "image"

var frameBounds = image.Rect(0, 0, colsPerFrame, rowsPerFrame)

// NewFrameImage returns a new image suitable for use with
// Lepton3.NextFrame().
func NewFrameImage() *image.Gray16 {
	return image.NewGray16(frameBounds)
}
