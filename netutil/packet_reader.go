package netutil

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/mkideal/log"
)

type PacketHandler func(b []byte)

// PacketReader reads a network packet
type PacketReader interface {
	Conn() net.Conn
	ReadPacket() (n int, err error)
	SetTimeout(d time.Duration)
}

type packetReader struct {
	conn          net.Conn
	timeout       time.Duration
	buf           []byte
	byte1         [LengthNeedSize]byte
	packetHandler PacketHandler
}

// NewPacketReader creates a PacketReader with net.Conn and PacketHandler
func NewPacketReader(conn net.Conn, packetHandler PacketHandler) PacketReader {
	return &packetReader{
		conn:          conn,
		packetHandler: packetHandler,
	}
}

func (r *packetReader) Conn() net.Conn             { return r.conn }
func (r *packetReader) SetTimeout(d time.Duration) { r.timeout = d }

const LengthNeedSize = 4

func (r *packetReader) ReadPacket() (int, error) {
	total := 0
	// set read timeout
	if r.timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.timeout))
	}
	// read packet length(lenof(packet.length)+lenof(packet.body)
	n, err := r.conn.Read(r.byte1[:])
	total += n
	if err != nil {
		log.Info("read error: %v", err)
		return total, err
	}
	length := binary.BigEndian.Uint32(r.byte1[:])
	log.Trace("packet length: %d", length)
	// read packet body
	if len(r.buf) < int(length) {
		r.buf = make([]byte, length)
	}
	n, err = r.conn.Read(r.buf)
	total += n
	if err != nil {
		log.Info("read error: %v", err)
		return total, err
	}
	log.Debug("read bytes number: %d", total)
	r.packetHandler(r.buf[:length])
	return total, nil
}
