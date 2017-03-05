package gateway

import (
	"errors"
	"net"
	"sync"

	"github.com/mkideal/pkg/netutil"
	"github.com/mkideal/pkg/netutil/protocol"
)

var (
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
)

type User interface {
	Id() int64
	GetSession() netutil.Session
	SetSession(netutil.Session)
	Authorized() bool
	SetAuthorized(bool)
	OnRecv([]byte)
	OnNewSession()
	OnQuitSession()
}

type Config struct {
	Addr         string
	Protocol     string // websocket or tcp
	Path         string // URL path for websocket protocol
	ConWriteSize int
}

type Gate struct {
	config  Config
	newUser func() User

	locker sync.RWMutex
	users  map[int64]User
}

func New(cfg Config, newUser func() User) *Gate {
	if cfg.Protocol == "" {
		cfg.Protocol = protocol.TCP
	}
	return &Gate{
		config:  cfg,
		users:   make(map[int64]User),
		newUser: newUser,
	}
}

func (gate *Gate) Run(async bool) error {
	var err error
	switch gate.config.Protocol {
	case protocol.TCP:
		err = netutil.ListenAndServeTCP(gate.config.Addr, gate.handleConn, async)
	case protocol.Websocket:
		err = netutil.ListenAndServeWebsocket(gate.config.Addr, gate.config.Path, gate.handleConn, async)
	default:
		err = ErrUnsupportedProtocol
	}
	return err
}

func (gate *Gate) Quit() {
	//TODO
}

func (gate *Gate) UserCount() int {
	gate.locker.RLock()
	defer gate.locker.RUnlock()
	return len(gate.users)
}

func (gate *Gate) handleConn(conn net.Conn) {
	id := conn.RemoteAddr().String()
	user := gate.newUser()
	reader := netutil.NewPacketReader(conn, user.OnRecv)
	session := netutil.NewRWSession(id, gate.config.ConWriteSize, reader)
	user.SetSession(session)
	session.Run(user.OnNewSession, user.OnQuitSession)
	user.SetSession(nil)
}

func (gate *Gate) AuthorizedUser(user User) {
	gate.locker.Lock()
	defer gate.locker.Unlock()
	id := user.Id()
	if oldUser, found := gate.users[id]; found {
		oldUser.SetAuthorized(false)
		delete(gate.users, id)
	}
	gate.users[id] = user
	user.SetAuthorized(true)
}

type UserVisitor func(User) (_break bool)

func (gate *Gate) ForEachUser(visitor UserVisitor) (_break bool) {
	gate.locker.RLock()
	defer gate.locker.RUnlock()
	for _, user := range gate.users {
		if _break = visitor(user); _break {
			break
		}
	}
	return
}

func (gate *Gate) Broadcast(receivers []int64, packet netutil.Packet) {
	gate.locker.RLock()
	defer gate.locker.RUnlock()
	if len(receivers) > 0 {
		for _, receiver := range receivers {
			if user, ok := gate.users[receiver]; ok {
				gate.Send(user, packet)
			}
		}
	} else {
		for _, user := range gate.users {
			gate.Send(user, packet)
		}
	}
}

func (gate *Gate) Send(user User, packet netutil.Packet) {
	if user.Authorized() {
		session := user.GetSession()
		if session != nil {
			session.Send(packet)
		}
	}
}
