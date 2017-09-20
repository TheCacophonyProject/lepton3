package main

import (
	"fmt"
	"time"

	"periph.io/x/periph/host"

	"github.com/TheCacophonyProject/lepton3"
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

	im := lepton3.NewFrameImage()
	t := time.Now()
	for i := 0; i < 90; i++ {
		fmt.Println(i)
		err := camera.NextFrame(im)
		checkErr("NextFrame", err)
	}
	fmt.Println(time.Since(t))

	// err = dumpHumanImage("lepton.png", im)
	// checkErr("dumpHumanImage", err)
}
