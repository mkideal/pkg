package smtpd

import (
	"bytes"
	"net"
	"net/textproto"
	"strings"

	"github.com/mkideal/pkg/debug"
)

var ServiceName = "Service ready"

// The NOOP, HELP, EXPN, VRFY, and RSET commands can be used at any time
// during a session, or without previously initializing a session.
const (
	NONE = ""

	EHLO = "EHLO"
	HELO = "HELO"
	MAIL = "MAIL"
	RCPT = "RCPT"
	DATA = "DATA"
	RSET = "RSET"
	VRFY = "VRFY"
	EXPN = "EXPN"
	HELP = "HELP"
	NOOP = "NOOP"
	QUIT = "QUIT"
	SIZE = "SIZE"
	AUTH = "AUTH"
)

type Server struct {
}

func NewServer() *Server {
	svr := new(Server)
	return svr
}

func (svr *Server) Start(addr string, listenCallback func(string)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if listenCallback != nil {
		listenCallback("listening on " + addr)
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		s := newSession(c)
		go s.run()
	}
}

// Sesion
type Session struct {
	nativeConn  net.Conn
	conn        *textproto.Conn
	prevCommand string
	didHello    bool
	buff        *bytes.Buffer
}

func newSession(conn net.Conn) *Session {
	s := new(Session)
	s.nativeConn = conn
	s.conn = textproto.NewConn(conn)
	s.buff = bytes.NewBufferString("")
	return s
}

func (s *Session) run() {
	s.conn.PrintfLine("%3d %s", CodeServiceReady, ServiceName)
	for {
		line, err := s.conn.ReadLine()
		if err != nil {
			s.conn.Close()
			return
		}
		cmd := ""
		args := ""
		strs := strings.SplitN(line, " ", 2)
		if len(strs) > 0 {
			cmd = strs[0]
		}
		if len(strs) > 1 {
			args = strs[1]
		}
		if quit := s.dispatch(cmd, args); quit {
			s.conn.Close()
			return
		}
	}
}

func (s *Session) dispatch(cmd, args string) (quit bool) {
	debug.Debugf("recv command: %q, args: %q", cmd, args)
	switch cmd {
	case NOOP:
		s.responseOK()
	case HELP:
		s.commandNotImplemented(cmd)
	case EXPN:
		s.commandNotImplemented(cmd)
	case VRFY:
		s.commandNotImplemented(cmd)
	case RSET:
		s.commandNotImplemented(cmd)

	case HELO:
		fallthrough
	case EHLO:
		s.onHello(args)

	case QUIT:
		quit = s.onQuit()

	default:
		s.commandNotImplemented(cmd)
	}
	return quit
}

func (s *Session) commandNotImplemented(cmd string) {
	s.conn.PrintfLine("502 command not implemented")
}

func (s *Session) responseOK() {
	s.conn.PrintfLine("250 OK")
}

func (s *Session) onHello(arg string) {
	s.didHello = true

	if len(arg) > 0 {
		s.responseOK()
	} else {
		s.conn.PrintfLine("500 syntax error")
	}
}

func (s *Session) onQuit() bool {
	s.conn.PrintfLine("221 bye")
	return true
}
