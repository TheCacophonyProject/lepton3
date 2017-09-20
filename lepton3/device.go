package lepton3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

const (
	vospiPacketSize  = 164
	packetsPerRead   = 200
	transferSize     = vospiPacketSize * packetsPerRead
	packetBufferSize = 1024

	packetHeaderDiscard = 0x0F00
	packetNumMask       = 0x0FFF
)

func New() *Dev {
	return new(Dev)
}

type Dev struct {
	spiPort  spi.PortCloser
	spiConn  spi.Conn
	packetCh chan []byte
	done     chan struct{}
	wg       sync.WaitGroup
}

func (d *Dev) ReadFrame() error {
	d.open()
	defer d.close()
	d.readFrame()
	return nil
}

func (d *Dev) open() error {
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

func (d *Dev) close() {
	d.stopStream()

	if d.spiPort != nil {
		d.spiPort.Close()
	}
	d.spiConn = nil
}

func (d *Dev) reset() error {
	fmt.Println("RESET")
	d.close()
	time.Sleep(200 * time.Millisecond)
	return d.open()
}

func (d *Dev) startStream() {
	d.packetCh = make(chan []byte, packetBufferSize)
	d.done = make(chan struct{})
	d.wg.Add(1)

	go func() {
		defer d.wg.Done()
		for {
			rx := make([]byte, transferSize) // XXX don't allocate each time - ring buffer
			if err := d.spiConn.Tx(nil, rx); err != nil {
				// XXX how to report back errors
				fmt.Printf("Tx failed: %v\n", err)
				return
			}
			for i := 0; i < len(rx); i += vospiPacketSize {
				select {
				case <-d.done:
					return
				case d.packetCh <- rx[i : i+vospiPacketSize]:
				}
			}
		}
	}()
}

func (d *Dev) stopStream() {
	// XXX don't call this if the stream goroutine isn't running
	close(d.done)
	d.wg.Wait()
}

func (d *Dev) readFrame() error {
	// XXX open must have been called first
	// XXX timeout when nothing valid for some time
	// XXX CRC checks

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

		segmentNum := int(packet[0] >> 4)
		complete, err := f.addPacket(packetNum, segmentNum)
		if err != nil {
			fmt.Printf("addPacket: %v\n", err)
			if err := d.reset(); err != nil {
				return err
			}
			f = newFrame()
		} else if complete {
			fmt.Printf("frame\n")
			return nil
		}
	}
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
	// XXX might not necessary with CRC check
	if packetNum == 0 && packet[2] == 0 && packet[3] == 0 {
		// fmt.Println("skipping 0 packet")
		return -1, nil
	}
	if packetNum > 60 {
		return -1, errors.New("invalid packet number")
	}
	return packetNum, nil
}

func newFrame() *frame {
	return &frame{
		packetNum: -1,
	}
}

type frame struct {
	packetNum  int
	segmentNum int
}

func (f *frame) addPacket(packetNum int, segmentNum int) (bool, error) {
	// fmt.Println(packetNum)
	if !f.isValidSeq(packetNum) {
		return false, fmt.Errorf("out of order packet: %d -> %d", f.packetNum, packetNum)
	}

	switch packetNum {
	case 20:
		f.segmentNum = segmentNum
	case 59:
		f.packetNum = -1
		if f.segmentNum == 4 {
			return true, nil
		}
	}
	f.packetNum = packetNum
	return false, nil
}

func (f *frame) isValidSeq(packetNum int) bool {
	if packetNum == 0 && f.packetNum == 59 {
		return true
	}
	return packetNum == f.packetNum+1
}
