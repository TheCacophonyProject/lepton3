package lepton3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

// XXX put code in the right place
// XXX use fixed errors where possible
// XXX deal with printfs (enable a debug mode or something?)
// XXX document copy minimisation
// XXX allow speed to be selected
// XXX profiling
// XXX vendoring

const (
	// Video Over SPI packets
	vospiHeaderSize = 4 // 2 byte ID, 2 byte CRC
	vospiDataSize   = 160
	vospiPacketSize = vospiHeaderSize + vospiDataSize

	// Packets, segments and frames
	colsPerFrame      = 160
	rowsPerFrame      = 120
	packetsPerSegment = 60
	segmentsPerFrame  = 4
	packetsPerFrame   = segmentsPerFrame * packetsPerSegment
	colsPerPacket     = colsPerFrame / 2

	// SPI transfer
	packetsPerRead   = 200 // XXX play around with this to check effect on CPU load and reliability
	transferSize     = vospiPacketSize * packetsPerRead
	packetBufferSize = 1024

	// Packet bitmasks
	packetHeaderDiscard = 0x0F00
	packetNumMask       = 0x0FFF
)

// New returns a new Lepton3 instance.
func New() *Lepton3 {
	return new(Lepton3)
}

// Lepton3 manages a connection to an FLIR Lepton 3 camera. It is not
// goroutine safe.
type Lepton3 struct {
	spiPort  spi.PortCloser
	spiConn  spi.Conn
	packetCh chan []byte
	done     chan struct{}
	wg       sync.WaitGroup
}

// Open initialises the SPI connection and starts streaming packets
// from the camera.
func (d *Lepton3) Open() error {
	spiPort, err := spireg.Open("")
	if err != nil {
		return err
	}
	spiConn, err := spiPort.Connect(30000000, spi.Mode3, 8)
	if err != nil {
		spiPort.Close()
		return err
	}

	d.spiPort = spiPort
	d.spiConn = spiConn

	d.startStream()
	return nil
}

// Close stops streaming of packets from the camera and closes the SPI
// device connection. It must only be called if streaming was started
// with Open().
func (d *Lepton3) Close() {
	d.stopStream()

	if d.spiPort != nil {
		d.spiPort.Close()
	}
	d.spiConn = nil
}

// NextFrame returns the next frame from the camera into the image
// provided.
//
// The output image is provided (rather than being created by
// NextFrame) to minimise memory allocations. Use NewFrameImage() to
// create an image suitable for use with NextFrame().
//
// NextFrame should only be called after a successful call to
// Open(). Although there is some internal buffering of camera
// packets, NextFrame must be called frequently enough to ensure
// frames are not lost.
func (d *Lepton3) NextFrame(im *image.Gray16) error {
	// XXX timeout when nothing valid for some time

	f := newFrame()
	for {
		packet := <-d.packetCh

		packetNum, err := validatePacket(packet)
		if err != nil {
			fmt.Println(err)
			if err := d.reset(); err != nil {
				return err
			}
			f = newFrame()
			continue
		} else if packetNum < 0 {
			continue
		}

		complete, err := f.nextPacket(packetNum, packet)
		if err != nil {
			fmt.Printf("addPacket: %v\n", err)
			if err := d.reset(); err != nil {
				return err
			}
			f = newFrame()
		} else if complete {
			f.writeImage(im)
			return nil
		}
	}
}

// Snapshot is convenience method for capturing a single frame. It
// should not be called if streaming is already active.
func (d *Lepton3) Snapshot() (*image.Gray16, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	defer d.Close()
	im := NewFrameImage()
	if err := d.NextFrame(im); err != nil {
		return nil, err
	}
	return im, nil
}

func (d *Lepton3) reset() error {
	fmt.Println("RESET")
	d.Close()
	time.Sleep(200 * time.Millisecond)
	return d.Open()
}

func (d *Lepton3) startStream() {
	// XXX check how long the channel ends up getting under normal use
	d.packetCh = make(chan []byte, packetBufferSize)
	d.done = make(chan struct{})
	d.wg.Add(1)

	go func() {
		defer d.wg.Done()
		for {
			// XXX don't allocate each time - ring buffer
			rx := make([]byte, transferSize)
			if err := d.spiConn.Tx(nil, rx); err != nil {
				// XXX report back errors
				fmt.Printf("Tx failed: %v\n", err)
				return
			}
			for i := 0; i < len(rx); i += vospiPacketSize {
				// XXX skip invalid packets here?
				select {
				case <-d.done:
					return
				case d.packetCh <- rx[i : i+vospiPacketSize]:
				}
			}
		}
	}()
}

func (d *Lepton3) stopStream() {
	close(d.done)
	d.wg.Wait()
}

func validatePacket(packet []byte) (int, error) {
	header := binary.BigEndian.Uint16(packet)
	if header&0x8000 == 0x8000 {
		return -1, errors.New("first bit set on header")
	}
	if header&packetHeaderDiscard == packetHeaderDiscard {
		return -1, nil
	}

	packetNum := int(header & packetNumMask)
	if packetNum > 60 {
		return -1, errors.New("invalid packet number")
	}

	// XXX might not necessary with CRC check
	if packetNum == 0 && packet[2] == 0 && packet[3] == 0 {
		return -1, nil
	}

	// XXX CRC checks

	return packetNum, nil
}

// XXX reuse a singe frame instead of recreating?
func newFrame() *frame {
	return &frame{
		packetNum:      -1,
		segmentNum:     0,
		segmentPackets: make([][]byte, packetsPerSegment),
		framePackets:   make([][]byte, packetsPerFrame),
	}
}

type frame struct {
	packetNum      int
	segmentNum     int
	segmentPackets [][]byte
	framePackets   [][]byte
}

func (f *frame) nextPacket(packetNum int, packet []byte) (bool, error) {
	if !f.sequential(packetNum) {
		return false, fmt.Errorf("out of order packet: %d -> %d", f.packetNum, packetNum)
	}

	// Store the packet data in current segment
	f.segmentPackets[packetNum] = packet[vospiHeaderSize:]

	switch packetNum {
	case 20:
		segmentNum := int(packet[0] >> 4)
		if segmentNum > 4 {
			return false, fmt.Errorf("invalid segment number: %d", segmentNum)
		}
		if segmentNum > 0 && segmentNum != f.segmentNum+1 {
			return false, fmt.Errorf("out of order segment")
		}
		f.segmentNum = segmentNum
	case 59:
		if f.segmentNum > 0 {
			// This should be fast as only slice headers for the
			// segment are being copied, not the packet data itself.
			copy(f.framePackets[(f.segmentNum-1)*packetsPerSegment:], f.segmentPackets)
		}
		if f.segmentNum == 4 {
			// Complete frame!
			return true, nil
		}
	}
	f.packetNum = packetNum
	return false, nil
}

func (f *frame) sequential(packetNum int) bool {
	if packetNum == 0 && f.packetNum == 59 {
		return true
	}
	return packetNum == f.packetNum+1
}

func (f *frame) writeImage(im *image.Gray16) {
	// XXX is there a faster way of doing this rather writing one
	// pixel at a time? Can we blat out a packet row at a time?
	// (remember endian-ness)
	for packetNum, packet := range f.framePackets {
		for i := 0; i < vospiDataSize; i += 2 {
			x := i >> 1 // divide 2
			if packetNum%2 == 1 {
				x += colsPerPacket
			}
			y := packetNum >> 1 // divide 2
			c := binary.BigEndian.Uint16(packet[i : i+2])
			im.SetGray16(x, y, color.Gray16{c})
		}
	}
}
