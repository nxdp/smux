// MIT License
//
// Copyright (c) 2016-2017 xtaci
//
// nxdp/smux: reduced header from 8 bytes to 5 bytes
//   - removed version byte (protocol upgrades are coordinated across controlled endpoints)
//   - StreamID shrunk from uint32 to uint16 (32767 initiated streams per side with the current odd/even allocator)
//
// Wire format:
//   xtaci/smux: | Ver(1) | Cmd(1) | Length(2) | StreamID(4) | = 8 bytes
//   nxdp/smux: | Cmd(1) | Length(2) | StreamID(2) |           = 5 bytes

package smux

import (
	"encoding/binary"
	"fmt"
)

const (
	cmdSYN byte = iota // stream open
	cmdFIN             // stream close, a.k.a EOF mark
	cmdPSH             // data push
	cmdNOP             // no operation
	cmdUPD             // notify bytes consumed by remote peer-end (flow control)
)

const (
	// data size of cmdUPD: |4B consumed| 4B window|
	szCmdUPD = 8
)

const (
	initialPeerWindow = 262144
)

const (
	// Cmd(1) + Length(2) + StreamID(2) = 5 bytes
	sizeOfCmd    = 1
	sizeOfLength = 2
	sizeOfSid    = 2
	headerSize   = sizeOfCmd + sizeOfLength + sizeOfSid
)

// Frame defines a packet to be multiplexed into a single connection
type Frame struct {
	cmd  byte
	sid  uint16
	data []byte
}

func newFrame(cmd byte, sid uint16) Frame {
	return Frame{cmd: cmd, sid: sid}
}

// rawHeader layout: [Cmd(1)][Length(2)][StreamID(2)]
type rawHeader [headerSize]byte

func (h rawHeader) Cmd() byte {
	return h[0]
}

func (h rawHeader) Length() uint16 {
	return binary.LittleEndian.Uint16(h[1:])
}

func (h rawHeader) StreamID() uint16 {
	return binary.LittleEndian.Uint16(h[3:])
}

func (h rawHeader) String() string {
	return fmt.Sprintf("Cmd:%d StreamID:%d Length:%d", h.Cmd(), h.StreamID(), h.Length())
}

type updHeader [szCmdUPD]byte

func (h updHeader) Consumed() uint32 {
	return binary.LittleEndian.Uint32(h[:])
}

func (h updHeader) Window() uint32 {
	return binary.LittleEndian.Uint32(h[4:])
}
