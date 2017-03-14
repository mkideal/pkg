package netutil

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/mkideal/log"
)

var (
	errLengthTooBig = errors.New("length too big")
)

type PacketHandler func(b []byte)

// PacketReader reads a network packet
type PacketReader interface {
	Conn() net.Conn
	ReadPacket() (n int, err error)
	SetTimeout(d time.Duration)
}

type packetReader struct {
	id            string
	conn          net.Conn
	timeout       time.Duration
	buf           []byte
	byte1         [LengthNeedSize]byte
	packetHandler PacketHandler
}

// NewPacketReader creates a PacketReader with net.Conn and PacketHandler
func NewPacketReader(conn net.Conn, packetHandler PacketHandler) PacketReader {
	return &packetReader{
		id:            conn.RemoteAddr().String(),
		conn:          conn,
		packetHandler: packetHandler,
	}
}

func (r *packetReader) Conn() net.Conn             { return r.conn }
func (r *packetReader) SetTimeout(d time.Duration) { r.timeout = d }

const (
	LengthNeedSize  = 4
	MaxPacketLength = 4 * 1024 * 1024 // 4M
)

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
		log.Info("%s: read error: %v", r.id, err)
		return total, err
	}
	length := binary.BigEndian.Uint32(r.byte1[:])
	log.Trace("%s: binary BigEndian decode %v as uint32 result: %d", r.id, r.byte1[:], length)
	if length > MaxPacketLength {
		log.Info("%s: length %d too big", r.id, length)
		return total, errLengthTooBig
	}
	// read packet body
	if len(r.buf) < int(length) {
		r.buf = make([]byte, length)
	}
	n, err = r.conn.Read(r.buf[:length])
	total += n
	if err != nil {
		log.Info("%s: read error: %v", r.id, err)
		return total, err
	}
	log.Debug("%s: read bytes number: %d", r.id, total)
	r.packetHandler(r.buf[:length])
	return total, nil
}
