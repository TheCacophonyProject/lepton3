package lepton3

import (
	"encoding/binary"

	"github.com/TheCacophonyProject/go-cptv/cptvframe"
)

type RawFrame [packetsPerFrame * vospiDataSize]byte

func (rf *RawFrame) FrameData() []byte {
	return rf[telemetryPacketCount*vospiDataSize:]
}

// ToFrame converts a RawFrame to a Frame.
func (rf *RawFrame) ToFrame(out *cptvframe.Frame) error {
	if err := ParseTelemetry(rf[:], &out.Status); err != nil {
		return err
	}

	rawPix := rf.FrameData()
	i := 0
	for y, row := range out.Pix {
		for x, _ := range row {
			out.Pix[y][x] = binary.BigEndian.Uint16(rawPix[i : i+2])
			i += 2
		}
	}

	return nil
}
