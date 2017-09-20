package main

import (
	"fmt"
	"time"

	"periph.io/x/periph/host"

	"golepton/lepton3"
)

func checkErr(label string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", err))
	}
}

// XXX implement snapshot and streaming options

func main() {
	_, err := host.Init()
	checkErr("host init", err)

	camera := lepton3.New()
	err = camera.Open()
	checkErr("Open", err)
	defer camera.Close()

	t := time.Now()
	for i := 0; i < 90; i++ {
		fmt.Println(i)
		_, err := camera.NextFrame()
		checkErr("NextFrame", err)
	}
	fmt.Println(time.Since(t))

	// err = dumpHumanImage("lepton.png", im)
	// checkErr("dumpHumanImage", err)
}
