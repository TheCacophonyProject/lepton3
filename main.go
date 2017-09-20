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

func main() {
	_, err := host.Init()
	checkErr("host init", err)

	dev := lepton3.New()
	err = dev.Open()
	checkErr("Open", err)
	defer dev.Close()

	t := time.Now()
	for i := 0; i < 90; i++ {
		fmt.Println(i)
		_, err := dev.NextFrame()
		checkErr("NextFrame", err)
	}
	fmt.Println(time.Since(t))

	// err = dumpHumanImage("lepton.png", im)
	// checkErr("dumpHumanImage", err)
}
