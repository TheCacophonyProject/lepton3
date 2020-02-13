package lepton3

import (
	"encoding/binary"

	"github.com/TheCacophonyProject/go-cptv/cptvframe"
)

// NewRawFrame returns a correctly sized byte slice for holding a
// single Lepton 3 frame.
func NewRawFrame() []byte {
	return make([]byte, BytesPerFrame)
}

// ParseRawFrame converts a byte slice containing a raw Lepton 3 frame
// into a cptvframe.Frame. The result is writing into the Frame
// provided.
func ParseRawFrame(raw []byte, out *cptvframe.Frame) error {
	if err := ParseTelemetry(raw, &out.Status); err != nil {
		return err
	}

	rawPix := raw[telemetryBytes:]
	i := 0
	for y, row := range out.Pix {
		for x := range row {
			out.Pix[y][x] = binary.BigEndian.Uint16(rawPix[i : i+2])
			i += 2
		}
	}

	return nil
}
