package main

import (
	"fmt"

	"periph.io/x/periph/host"

	"golepton/lepton3"
)

func checkErr(label string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", err))
	}
}

func main() {
	_, err := host.Init()
	checkErr("host init", err)

	dev := lepton3.New()
	// im, err := dev.ReadFrame()
	// checkErr("ReadFrame", err)

	i := 0
	imCh, _, err := dev.StreamFrames()
	checkErr("StreamFrames", err)
	for {
		select {
		case im := <-imCh:
			var _ = im
			// dumpHumanImage(fmt.Sprintf("%04d.png", i), im)
			// checkErr("dumpHumanImage", err)
			i++
		}
	}

	// err = dumpHumanImage("lepton.png", im)
	// checkErr("dumpHumanImage", err)
}
