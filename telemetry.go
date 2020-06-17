// Copyright 2018 The Cacophony Project. All rights reserved.
// Use of this source code is governed by the Apache License Version 2.0;
// see the LICENSE file for further details.

// Some source code in this file comes from the periph project
// (https://periph.io/).

package lepton3

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/TheCacophonyProject/go-cptv/cptvframe"
)

// These are the valid values for the Telemetry.FFCState field.
const (
	FFCNever    = "never"
	FFCImminent = "imminent"
	FFCRunning  = "running"
	FFCComplete = "complete"
)

// ParseTelemetry converts a slice containing raw Lepton 3 telemetry
// data into a Telemetry struct.
func ParseTelemetry(raw []byte, t *cptvframe.Telemetry) error {
	var tw telemetryWords
	if err := binary.Read(bytes.NewBuffer(raw), Big16, &tw); err != nil {
		return err
	}

	t.TimeOn = tw.TimeOn.ToD()
	t.FFCState = statusToFFCState(tw.StatusBits)
	t.FrameCount = int(tw.FrameCounter)
	t.FrameMean = tw.FrameMean
	t.TempC = tw.FPATemp.ToC()
	t.LastFFCTempC = tw.FPATempLastFFC.ToC()
	t.LastFFCTime = tw.TimeCounterLastFFC.ToD()
	return nil
}

type LeptonSoftwareRevision struct {
	Gpp_major uint8
	Gpp_minor uint8
	Gpp_build uint8
	Dsp_major uint8
	Dsp_minor uint8
	Dsp_build uint8
	Reserved  [2]uint8
}

type telemetryWords struct {
	TelemetryRevision  uint16                 // 0  *
	TimeOn             durationMS             // 1  *
	StatusBits         uint32                 // 3  * Bit field
	Reserved5          [8]uint16              // 5  *
	SoftwareRevision   LeptonSoftwareRevision // 13 - Junk.
	Reserved17         [3]uint16              // 17 *
	FrameCounter       uint32                 // 20 *
	FrameMean          uint16                 // 22 * The average value from the whole frame
	FPATempCounts      uint16                 // 23
	FPATemp            centiK                 // 24 *
	HousingTempRaw     uint16
	HousingTemp        centiK
	Reserved25         [2]uint16  // 25
	FPATempLastFFC     centiK     // 29
	TimeCounterLastFFC durationMS // 30 *
	HousingTempLastFFC centiK
}

// durationMS is duration in millisecond.
//
// It is an implementation detail of the protocol.
type durationMS uint32

// ToD converts a millisecond based timing to time.Duration.
func (d durationMS) ToD() time.Duration {
	return time.Duration(d) * time.Millisecond
}

// centiK is temperature in 0.01°K
//
// It is an implementation detail of the protocol.
type centiK uint16

// ToC converts a Kelvin measurement to Celsius.
func (c centiK) ToC() float64 {
	return float64(int(c)-27315) / 100
}

const statusFFCStateMask uint32 = 3 << 4
const statusFFCStateShift uint32 = 4

func statusToFFCState(status uint32) string {
	bits := status & statusFFCStateMask >> statusFFCStateShift
	switch bits {
	case 0:
		return FFCNever
	case 1:
		return FFCImminent
	case 2:
		return FFCRunning
	default:
		return FFCComplete
	}
}
