package poker

import (
	"fmt"
	"sync"

	"github.com/mkideal/pkg/game"
	"github.com/mkideal/pkg/math/random"
)

type Landlord3Player interface {
	// embed game.Player
	game.Player
	// Pos returns position of player:0,1,2
	Pos() int
	// IsLandlord reports whther player is landlord
	IsLandlord() bool
	SetIsLandlord(yes bool)
	// Pokers returns current all pokers of player
	Pokers() []Poker
	// AddPokers adds pokers
	AddPokers(pokers []Poker)
	// Multiple returns current multiple
	Multiple() int

	Double(multiple int)
}

// landlord for three players
type Landlord3 struct {
	uuid         string
	players      []Landlord3Player
	lastCards    []Poker
	randomSource random.Source
}

// NewLandlord3 new Landlord3, panic if len(players) != 3
func NewLandlord3(uuid string, players ...Landlord3Player) *Landlord3 {
	f := &Landlord3{
		uuid:    uuid,
		players: players,
	}
	if len(players) != f.NumPlayer() {
		panic(fmt.Sprintf("length of players is not %d", f.NumPlayer()))
	}
	return f
}

func (f Landlord3) NumPlayer() int { return 3 }

// UUID returns uuid of fighting
func (f Landlord3) UUID() string { return f.uuid }

// SetRandomSource sets random source
func (f *Landlord3) SetRandomSource(source random.Source) { f.randomSource = source }

var landlord3Nums = []int{17, 17, 17}

// Start starts bureau
func (f *Landlord3) Start() error {
	// deal pokers
	pokers, remains := Deal(landlord3Nums, GetPokers(), f.randomSource)
	for i, player := range f.players {
		player.Reset()
		player.AddPokers(pokers[i])
	}
	f.lastCards = remains

	// randomly chose landlord
	pos := random.Intn(f.NumPlayer(), f.randomSource)
	_ = pos

	return nil
}

func (f *Landlord3) Shutdown(g *sync.WaitGroup) {
	g.Done()
}

func (f *Landlord3) Recv(cmd game.Command) {
	f.onCommand(cmd)
}

func (f *Landlord3) Double(multiple int) {
	for _, player := range f.players {
		player.Double(multiple)
	}
}

func (f Landlord3) gameover() bool {
	for _, player := range f.players {
		if len(player.Pokers()) == 0 {
			return true
		}
	}
	return false
}

func (f *Landlord3) onGameover(force bool) {
	//TODO
}

func (f *Landlord3) onCommand(cmd game.Command) {
	if cmd.Pos >= 0 && cmd.Pos < len(f.players) {
		player := f.players[cmd.Pos]
		_ = player
		//TODO
	} else {
		//TODO
	}
}
