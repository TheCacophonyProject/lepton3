// Copyright 2020 The Cacophony Project. All rights reserved.
// Use of this source code is governed by the Apache License Version 2.0;
// see the LICENSE file for further details.

// This file mirrors some types and consts from periph.io.

package lepton3

import (
	"periph.io/x/periph/devices"
	"periph.io/x/periph/devices/lepton/cci"
)

type FFCShutterMode = cci.FFCShutterMode

const (
	FFCShutterModeManual   = cci.FFCShutterModeManual
	FFCShutterModeAuto     = cci.FFCShutterModeAuto
	FFCShutterModeExternal = cci.FFCShutterModeExternal
)

type FFCMode = cci.FFCMode

const (
	ShutterTempLockoutStateInactive = cci.ShutterTempLockoutStateInactive
	ShutterTempLockoutStateHigh     = cci.ShutterTempLockoutStateHigh
	ShutterTempLockoutStateLow      = cci.ShutterTempLockoutStateLow
)

type Celsius = devices.Celsius

// CelsiusFromFloat creates a new Celsius from a floating point
// value. This is used for temperature fields in FFCMode.
func CelsiusFromFloat(c float64) Celsius {
	return Celsius(int(c * 1000))
}
