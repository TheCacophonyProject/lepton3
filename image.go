// Copyright 2017 The Cacophony Project. All rights reserved.
// Use of this source code is governed by the Apache License Version 2.0;
// see the LICENSE file for further details.

package lepton3

import "image"

var frameBounds = image.Rect(0, 0, colsPerFrame, rowsPerFrame)

// NewFrameImage returns a new image suitable for use with
// Lepton3.NextFrame().
func NewFrameImage() *image.Gray16 {
	return image.NewGray16(frameBounds)
}
