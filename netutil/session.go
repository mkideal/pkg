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
	Send(Packet)
	Run(onNewSession, onQuitSession func())
	Quit()
}

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
func (ws *WSession) setClosed()      { atomic.StoreInt32(&ws.closed, 1) }
func (ws *WSession) getClosed() bool { return atomic.LoadInt32(&ws.closed) == 1 }

func (ws *WSession) Send(p Packet) {
	if p.Len() > 0 && !ws.getClosed() {
		ws.writeChan <- p
	}
}

func (ws *WSession) write(b []byte) (err error) {
	/*if encrypter := ws.encrypter; encrypter != nil {
		start := 0
		length := 0
		for len(src[start:]) > 0 {
			length = len(src[start:])
			if len(ws.encryptBuf) < length {
				length = len(ws.encryptBuf)
			}
			encrypter.XORStream(ws.encryptBuf[:length], src[start:start+length])
			_, err = ws.conn.Write(ws.encryptBuf[:length])
			start += length
		}
	} else {*/
	_, err = ws.conn.Write(b)
	//}
	//b.Done()
	return
}

func (ws *WSession) startWriteLoop(startWrite, endWrite chan<- struct{}) {
	startWrite <- struct{}{}
	remain := 0
	for {
		if ws.getClosed() {
			remain = len(ws.writeChan)
			break
		}
		select {
		case p := <-ws.writeChan:
			err := ws.write(p.Bytes())
			if err != nil {
				ws.setClosed()
			}
		case <-time.After(time.Second):
		}
	}

	for i := 0; i < remain; i++ {
		p := <-ws.writeChan
		err := ws.write(p.Bytes())
		if err != nil {
			break
		}
	}

	ws.conn.Close()
	log.Debug("WSession startWriteLoop end")
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
	startRead <- struct{}{}
	for {
		_, err := s.packetReader.ReadPacket()
		if err != nil {
			s.setClosed()
		}
		if s.getClosed() {
			break
		}
	}
	log.Debug("RWSession startReadLoop end")
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
