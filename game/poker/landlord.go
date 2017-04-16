package poker

import (
	"fmt"
	"sync/atomic"

	"github.com/mkideal/pkg/game"
	"github.com/mkideal/pkg/math/random"
)

// see: https://zh.wikipedia.org/wiki/鬥地主
// or
// see: https://en.wikipedia.org/wiki/Dou_dizhu

// landlord                       : 地主
// peasants                       : 农民
// bid landlord                   : 叫地主
// shuffle cards                  : 洗牌
// deal cards                     : 发牌
// leftover wild cards            : 底牌
// rocket                         : 王炸(火箭)
// bomb                           : 炸弹
// solo                           : 单牌
// pair                           : 对子
// trio                           : 三张
// four                           : 四张
// chain                          : 顺子
// pairs chain                    : 连对
// trio with single card          : 三带一
// trio with pair                 : 三带二
// airplain                       : 飞机
// airplain with small wings      : 飞机带小翼
// airplain with large wings      : 飞机带大翼
// four with two single cards     : 四带两张
// four with two pairs            : 四带两对
// space shuttle                  : 航天飞机
// space shuttle with small wings : 航天飞机带小翼
// space shuttle with large wings : 航天飞机带小翼
// spring                         : 春天

type Landlord3Player interface {
	// embed game.Player
	game.Player
	// Pos returns position of player:0,1,2
	Pos() int
	// SetContext sets Landlord3 as context
	SetContext(landlord3 *Landlord3)
	// IsLandlord reports whther player is landlord
	IsLandlord() bool
	SetIsLandlord(yes bool)
	// Pokers returns current all pokers of player
	Pokers() []Poker
	// AddPokers adds pokers
	AddPokers(pokers []Poker)
	// Multiple returns current multiple
	Multiple() int
	// OnGameover handle result
	OnGameover(winnerPos int) error

	// commands

	Double(multiple int) error
	BidLandlord(score int) error
}

type LandlordStage int

const (
	LandlordReady LandlordStage = iota
	LandlordBid
	LandlordDouble
	LandlordPlay
)

// landlord for three players
type Landlord3 struct {
	uuid         string
	players      []Landlord3Player
	lastCards    []Poker
	randomSource random.Source

	playing  int32
	stage    LandlordStage
	lastTurn int
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

// NumPlayer returns number of players
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
		player.SetContext(f)
		player.AddPokers(pokers[i])
	}
	f.lastCards = remains

	// randomly choose landlord
	pos := random.Intn(f.NumPlayer(), f.randomSource)

	f.stage = LandlordReady
	f.lastTurn = pos
	f.setPlaying(true)
	return nil
}

func (f *Landlord3) setPlaying(on bool) bool {
	if on {
		return atomic.CompareAndSwapInt32(&f.playing, 0, 1)
	} else {
		return atomic.CompareAndSwapInt32(&f.playing, 1, 0)
	}
}

func (f *Landlord3) Shutdown() error {
	return f.beforeGameover(true)
}

// Double doubles all players with multiple `multiple`
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

func (f *Landlord3) beforeGameover(force bool) error {
	if !f.setPlaying(false) {
		return game.ErrGameover
	}
	if !force {
		winnerPos := -1
		for i := range f.players {
			if len(f.players[i].Pokers()) == 0 {
				winnerPos = i
			}
		}
		if winnerPos == -1 {
			return game.ErrUnexpectedGameover
		}
		for i := range f.players {
			f.players[i].OnGameover(winnerPos)
		}
	}
	return nil
}

// GetPlayer gets player by pos
func (f *Landlord3) GetPlayer(pos int) Landlord3Player {
	n := f.NumPlayer()
	pos = pos % n
	if pos < 0 {
		pos += n
	}
	return f.players[pos]
}

// GetPrevPlayer gets previous player by current player pos
func (f *Landlord3) GetPrevPlayer(pos int) Landlord3Player {
	return f.GetPlayer(pos - 1)
}

// GetNextPlayer gets next player by current player pos
func (f *Landlord3) GetNextPlayer(pos int) Landlord3Player {
	return f.GetPlayer(pos + 1)
}

// receive command

func (f *Landlord3) OnBidLandlord(pos int, score int) error {
	player := f.GetPlayer(pos)
	return player.BidLandlord(score)
}

func (f *Landlord3) OnDouble(pos int, multiple int) error {
	player := f.GetPlayer(pos)
	return player.Double(multiple)
}
