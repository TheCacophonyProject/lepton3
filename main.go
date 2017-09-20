package main

import (
	"fmt"

	"golepton/lepton3"

	"periph.io/x/periph/host"
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
	dev.ReadFrame()
}
