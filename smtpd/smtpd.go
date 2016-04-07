package smtpd

import (
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strings"

	"github.com/mkideal/pkg/debug"
)

var ServiceInfo = "Service ready"

var (
	fromRegexp = regexp.MustCompile("[Ff][Rr][Oo][Mm]:(.+)")
	toRegexp   = regexp.MustCompile("[Tt][Oo]:(.+)")
)

// The NOOP, HELP, EXPN, VRFY, and RSET commands can be used at any time
// during a session, or without previously initializing a session.
const (
	NONE = ""

	EHLO     = "EHLO"
	HELO     = "HELO"
	MAIL     = "MAIL"
	RCPT     = "RCPT"
	DATA     = "DATA"
	RSET     = "RSET"
	VRFY     = "VRFY"
	EXPN     = "EXPN"
	HELP     = "HELP"
	NOOP     = "NOOP"
	QUIT     = "QUIT"
	SIZE     = "SIZE"
	AUTH     = "AUTH"
	STARTTLS = "STARTTLS"
)

type Repository interface {
	SaveEmail(from, tos string, data []byte) error
}

type Server struct {
	repo Repository
}

func NewServer(repo Repository) *Server {
	svr := new(Server)
	svr.repo = repo
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
		s := newSession(svr, c)
		go s.run()
	}
}

// Sesion
type Session struct {
	svr        *Server
	nativeConn net.Conn
	conn       *textproto.Conn

	// whether the Session is using TLS
	tls bool

	// supported extensions
	ext map[string]string

	didHello bool

	starttls []byte
	auth     []byte
	mail     []byte
	rcpt     []byte
	data     []byte
}

func newSession(svr *Server, conn net.Conn) *Session {
	s := new(Session)
	s.svr = svr
	s.nativeConn = conn
	s.conn = textproto.NewConn(conn)
	return s
}

func (s *Session) run() {
	s.conn.PrintfLine("%3d %s", CodeServiceReady, ServiceInfo)
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
		s.onNoOp()
	case HELP:
		s.commandNotImplemented(cmd)
	case EXPN:
		s.commandNotImplemented(cmd)
	case VRFY:
		s.commandNotImplemented(cmd)
	case RSET:
		s.onRset()

	case HELO:
		s.onHelo(args)
	case EHLO:
		s.onEhlo(args)

	case QUIT:
		quit = s.onQuit()

	case STARTTLS:
		quit = s.onStartTLS(args)

	case AUTH:
		quit = s.onAuth(args)

	case MAIL:
		quit = s.onMail(args)

	case RCPT:
		quit = s.onRcpt(args)

	case DATA:
		quit = s.onData(args)

	case NONE:
		// do nothing

	default:
		s.commandNotImplemented(cmd)
	}
	return quit
}

func (s *Session) commandNotImplemented(cmd string) {
	s.conn.PrintfLine("502 command not implemented")
}

func (s *Session) complete() (quit bool) {
	quit = true
	return
}

//------------------
// command handlers
//------------------

// HELO
func (s *Session) onHelo(args string) {
	s.didHello = true

	if len(args) > 0 {
		s.responseOK()
	} else {
		s.responseSyntaxError()
	}
}

// EHLO
func (s *Session) onEhlo(args string) {
	s.onHelo(args)
}

// NOOP
func (s *Session) onNoOp() {
	s.responseOK()
}

// RSET
func (s *Session) onRset() {
	s.mail = nil
	s.rcpt = nil
	s.auth = nil
	s.starttls = nil
	s.data = nil
}

// STARTTLS
func (s *Session) onStartTLS(args string) (quit bool) {
	s.commandNotImplemented(STARTTLS)
	return
}

// AUTH
func (s *Session) onAuth(args string) (quit bool) {
	s.commandNotImplemented(AUTH)
	return
}

// MAIL
func (s *Session) onMail(args string) (quit bool) {
	matchResult := fromRegexp.FindStringSubmatch(args)
	if matchResult == nil || len(matchResult) != 2 {
		s.responseSyntaxError()
		return
	}
	s.mail = []byte(args)
	return
}

// RCPT
func (s *Session) onRcpt(args string) (quit bool) {
	matchResult := toRegexp.FindStringSubmatch(args)
	if matchResult == nil || len(matchResult) != 2 {
		s.responseSyntaxError()
		return
	}
	s.rcpt = []byte(args)
	return
}

// DATA
func (s *Session) onData(args string) (quit bool) {
	if args == "" {
		s.responseStartMailInput()
		return
	}
	if args == "." {
		return s.complete()
	}
	if s.data == nil {
		s.data = []byte(args)
	} else {
		s.data = append(s.data, []byte(args)...)
	}
	return
}

// QUIT
func (s *Session) onQuit() bool {
	s.responseQuit()
	return true
}

//----------
// response
//----------

func (s *Session) responseOK() {
	s.printf("%3d OK", CodeOK)
}

func (s *Session) responseSyntaxError() {
	s.printf("%3d syntax error", CodeSyntaxError)
}

func (s *Session) responseQuit() {
	s.printf("%3d bye", CodeServiceClosing)
}

func (s *Session) responseStartMailInput() {
	s.printf("%3d start mail input", CodeStartMailInput)
}

func (s *Session) printf(format string, args ...interface{}) {
	resp := fmt.Sprintf(format, args...)
	debug.Debugf("resp: %s", resp)
	s.conn.PrintfLine(resp)
}
