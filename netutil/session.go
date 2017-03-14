package netutil

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/mkideal/log"
)

type Packet interface {
	Len() int
	Bytes() []byte
}

type BytesPacket []byte

func (p BytesPacket) Len() int      { return len(p) }
func (p BytesPacket) Bytes() []byte { return []byte(p) }

type Session interface {
	Id() string
	Closed() bool
	Send(Packet)
	Run(onNewSession, onQuitSession func())
	Quit()
}

type nullSession struct{}

var NullSession = nullSession{}

func (session nullSession) Id() string         { return "" }
func (session nullSession) Send(Packet)        {}
func (session nullSession) Run(func(), func()) {}
func (session nullSession) Quit()              {}

// Write-only Session
type WSession struct {
	conn   net.Conn
	id     string
	closed int32

	writeQuit chan struct{}
	writeChan chan Packet
}

func NewWSession(id string, conn net.Conn, conWriteSize int) *WSession {
	if conWriteSize <= 0 {
		conWriteSize = 64
	}
	return &WSession{
		conn:      conn,
		id:        id,
		writeQuit: make(chan struct{}),
		writeChan: make(chan Packet, conWriteSize),
	}
}

func (ws *WSession) Id() string      { return ws.id }
func (ws *WSession) Closed() bool    { return ws.getClosed() }
func (ws *WSession) setClosed()      { atomic.StoreInt32(&ws.closed, 1) }
func (ws *WSession) getClosed() bool { return atomic.LoadInt32(&ws.closed) == 1 }

func (ws *WSession) Send(p Packet) {
	if p.Len() > 0 && !ws.getClosed() {
		ws.writeChan <- p
	}
}

func (ws *WSession) startWriteLoop(startWrite, endWrite chan<- struct{}) {
	startWrite <- struct{}{}
	remain := 0
	id := ws.Id()
	log.Debug("session %s write loop begin", id)
	for {
		if ws.getClosed() {
			remain = len(ws.writeChan)
			break
		}
		select {
		case p := <-ws.writeChan:
			_, err := ws.conn.Write(p.Bytes())
			if err != nil {
				ws.setClosed()
			}
		case <-time.After(time.Second):
		}
	}

	for i := 0; i < remain; i++ {
		p := <-ws.writeChan
		_, err := ws.conn.Write(p.Bytes())
		if err != nil {
			break
		}
	}

	ws.conn.Close()
	log.Debug("session %s write loop end", id)
	endWrite <- struct{}{}
}

func (ws *WSession) Run(onNewSession, onQuitSession func()) {
	startWrite := make(chan struct{})
	endWrite := make(chan struct{})

	go ws.startWriteLoop(startWrite, endWrite)
	<-startWrite

	if onNewSession != nil {
		onNewSession()
	}

	<-endWrite

	if ws.conn != nil {
		ws.conn.Close()
	}

	if onQuitSession != nil {
		onQuitSession()
	}
}

func (ws *WSession) Quit() {
	ws.setClosed()
}

// Readable and Writable Session
type RWSession struct {
	*WSession
	packetReader PacketReader
}

func NewRWSession(
	id string,
	conWriteSize int,
	packetReader PacketReader,
) *RWSession {
	s := new(RWSession)
	conn := packetReader.Conn()
	s.WSession = NewWSession(id, conn, conWriteSize)
	s.packetReader = packetReader
	return s
}

func (s *RWSession) startReadLoop(startRead, endRead chan<- struct{}) {
	id := s.Id()
	startRead <- struct{}{}
	log.Debug("session %s read loop begin", id)
	for {
		_, err := s.packetReader.ReadPacket()
		if err != nil {
			s.setClosed()
		}
		if s.getClosed() {
			break
		}
	}
	log.Debug("session %s read loop end", id)
	endRead <- struct{}{}
}

func (s *RWSession) Run(onNewSession, onQuitSession func()) {
	startRead := make(chan struct{})
	startWrite := make(chan struct{})
	endRead := make(chan struct{})
	endWrite := make(chan struct{})

	go s.startReadLoop(startRead, endRead)
	go s.startWriteLoop(startWrite, endWrite)

	<-startRead
	<-startWrite

	if onNewSession != nil {
		onNewSession()
	}

	<-endRead
	<-endWrite

	if s.conn != nil {
		s.conn.Close()
	}

	if onQuitSession != nil {
		onQuitSession()
	}
}
