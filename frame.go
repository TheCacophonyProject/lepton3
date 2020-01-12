// Copyright 2018 The Cacophony Project. All rights reserved.
// Use of this source code is governed by the Apache License Version 2.0;
// see the LICENSE file for further details.

package lepton3

import (
	"encoding/binary"
)

// RawFrame hold the raw bytes for single frame. This is helpful for
// transferring frames around. It can be converted to the more useful
// Frame.
type RawFrame [packetsPerFrame * vospiDataSize]byte

// Frame represents the thermal readings for a single frame.
type Frame struct {
	Pix    [][]uint16
	Status Telemetry
}

type CameraResolution interface {
	XRes() int
	YRes() int
	FPS() int
}

func NewFrame(c CameraResolution) *Frame {
	frame := new(Frame)
	frame.Pix = make([][]uint16, c.YRes())
	for i := range frame.Pix {
		frame.Pix[i] = make([]uint16, c.XRes())
	}
	return frame
}

// ToFrame converts a RawFrame to a Frame.
func (rf *RawFrame) ToFrame(out *Frame) error {
	if err := ParseTelemetry(rf[:], &out.Status); err != nil {
		return err
	}

	rawPix := rf[telemetryPacketCount*vospiDataSize:]
	i := 0
	for y, row := range out.Pix {
		for x, _ := range row {
			out.Pix[y][x] = binary.BigEndian.Uint16(rawPix[i : i+2])
			i += 2
		}
	}

	return nil
}

// Copy sets current frame as other frame
func (fr *Frame) Copy(orig *Frame) {
	fr.Status = orig.Status
	for y, row := range orig.Pix {
		copy(fr.Pix[y][:], row)
	}
}
