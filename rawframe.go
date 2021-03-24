package lepton3

import (
	"encoding/binary"
	"fmt"

	"github.com/TheCacophonyProject/go-cptv/cptvframe"
)

type BadFrameErr struct {
	Cause error
}

func (e *BadFrameErr) Error() string {
	return e.Cause.Error()
}

// NewRawFrame returns a correctly sized byte slice for holding a
// single Lepton 3 frame.
func NewRawFrame() []byte {
	return make([]byte, BytesPerFrame)
}

// ParseRawFrame converts a byte slice containing a raw Lepton 3 frame
// into a cptvframe.Frame. The result is writing into the Frame
// provided.
func ParseRawFrame(raw []byte, out *cptvframe.Frame, edgePixels int) error {
	if err := ParseTelemetry(raw, &out.Status); err != nil {
		return err
	}

	rawPix := raw[telemetryBytes:]
	i := 0

	for y, row := range out.Pix {
		for x := range row {
			out.Pix[y][x] = binary.BigEndian.Uint16(rawPix[i : i+2])
			onEdge := y < edgePixels || x < edgePixels || y >= (len(out.Pix)-edgePixels) || x >= (len(row)-edgePixels)
			if !onEdge && out.Pix[y][x] == 0 {
				err := fmt.Errorf("Bad pixel (%d,%d) of %d", y, x, out.Pix[y][x])
				return &BadFrameErr{err}
			}
			i += 2
		}
	}

	return nil
}
