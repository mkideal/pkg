package landlord

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/mkideal/pkg/game"
	"github.com/mkideal/pkg/game/poker"
	"github.com/mkideal/pkg/math/random"
)

var (
	ErrEmptyPokers     = errors.New("empty pokers")
	ErrInvalidCategory = errors.New("invalid category")
)

// Landlord3Player represents a player for playing landlord
type Landlord3Player interface {
	// embed game.Player
	game.Player
	// Pos returns position of player
	Pos() int
	// SetPos sets position of player
	SetPos(pos int)
	// SetContext sets Landlord3 as context
	SetContext(landlord3 *Landlord3)
	// OnGameover handle result
	OnGameover(winnerPos int) error
	// Double doubles score
	Double(state *Landlord3PlayerState, multiple int) error
	// BidLandlord bids landlord
	BidLandlord(state *Landlord3PlayerState, score int) (over bool, err error)
	// RevealCards reveals cards
	RevealCards(state *Landlord3PlayerState) error
	// PlayCards plays cards
	PlayCards(state *Landlord3PlayerState, pokers []uint8) error
}

// Landlord3Selector represents a selector for selecting landlord
type Landlord3Selector interface {
	// Init returns first landlord position
	Init(players []Landlord3Player) int
	// Over reports whther landlord has been selected
	Over(players []Landlord3Player) bool
	// Result returns position of landlord
	Result(players []Landlord3Player) int
}

// landlord for three players
type Landlord3 struct {
	randomSource random.Source
	players      []Landlord3Player
	state        Landlord3State
}

// NewLandlord3 new Landlord3, panic if len(players) != 3
func NewLandlord3(uuid string, players ...Landlord3Player) *Landlord3 {
	f := &Landlord3{players: players}
	f.state.Uuid = uuid
	if len(players) != f.NumPlayer() {
		panic(fmt.Sprintf("number of players MUST be %d, but got %d", f.NumPlayer(), len(players)))
	}
	f.state.Players = make([]Landlord3PlayerState, len(players))
	for i := range players {
		player := &f.state.Players[i]
		player.Id = players[i].ID()
		player.Pos = int8(i)
		player.Multiple = 1
	}
	return f
}

// NumPlayer returns number of players
func (f Landlord3) NumPlayer() int { return 3 }

// UUID returns uuid
func (f Landlord3) UUID() string { return f.state.Uuid }

// SetRandomSource sets random source
func (f *Landlord3) SetRandomSource(source random.Source) { f.randomSource = source }

var landlord3Nums = []int{17, 17, 17}

// Start starts bureau
func (f *Landlord3) Start() error {
	// deal pokers
	pokers, leftover := poker.Deal(landlord3Nums, poker.GetPokers(), f.randomSource)
	for i, player := range f.players {
		player.SetPos(i)
		player.Reset()
		player.SetContext(f)
		f.addPokers(player, &f.state.Players[i], pokers[i])
	}
	f.state.LeftoverWildCards = leftover

	// first landlord position chosen randomly
	pos := random.Intn(f.NumPlayer(), f.randomSource)
	f.state.FirstLandlordPos = int8(pos)

	f.state.Stage = Stage_Ready
	f.state.Turn = int8(pos)
	f.setPlaying(true)

	f.state.Stage = Stage_Bid
	return nil
}

func (f *Landlord3) addPokers(player Landlord3Player, state *Landlord3PlayerState, pokers []uint8) {
	state.Pokers = append(state.Pokers, pokers...)
}

func (f *Landlord3) setPlaying(on bool) bool {
	if on {
		return atomic.CompareAndSwapInt32(&f.state.Playing, 0, 1)
	} else {
		return atomic.CompareAndSwapInt32(&f.state.Playing, 1, 0)
	}
}

func (f *Landlord3) Shutdown() error {
	return f.beforeGameover(true)
}

// double doubles all players with multiple `multiple`
func (f *Landlord3) double(multiple int) {
	for i := range f.state.Players {
		f.state.Players[i].Multiple *= int32(multiple)
	}
}

func (f Landlord3) gameover() bool {
	for _, player := range f.state.Players {
		if len(player.Pokers) == 0 {
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
			if len(f.state.Players[i].Pokers) == 0 {
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

func (f *Landlord3) expectStages(stages ...Stage) error {
	for _, want := range stages {
		if want == f.state.Stage {
			return nil
		}
	}
	return game.ErrState
}

func (f *Landlord3) expectTurn(pos int) error {
	if f.state.Turn != int8(pos) {
		return game.ErrTurn
	}
	return nil
}

func (f *Landlord3) nextTurn() int {
	f.state.Turn = (f.state.Turn + 1) % int8(f.NumPlayer())
	return int(f.state.Turn)
}

func (f *Landlord3) selectLandlord() {
	//TODO
	f.state.Landlord = f.state.Turn
	f.state.Stage = Stage_Double
}

func (f *Landlord3) doubled(pos int) bool {
	for _, d := range f.state.Doubled {
		if d == int8(pos) {
			return true
		}
	}
	return false
}

func (f *Landlord3) storeLastPokers(pokers []uint8) {
	size := len(pokers)
	if len(f.state.LastHands.LastPokers) < size {
		f.state.LastHands.LastPokers = make([]uint8, size)
	}
	f.state.LastHands.LastPokers = f.state.LastHands.LastPokers[:size]
	copy(f.state.LastHands.LastPokers, pokers)
}

// receive command

func (f *Landlord3) OnBidLandlord(pos int, score int) error {
	player := f.GetPlayer(pos)
	// check stage
	if err := f.expectStages(Stage_Bid); err != nil {
		return err
	}
	// check turn
	if err := f.expectTurn(pos); err != nil {
		return err
	}
	// bid landlord
	if over, err := player.BidLandlord(&f.state.Players[player.Pos()], score); err != nil {
		return err
	} else if !over {
		f.nextTurn()
	} else {
		f.selectLandlord()
	}
	return nil
}

func (f *Landlord3) OnDouble(pos int, multiple int) error {
	// check stage
	if err := f.expectStages(Stage_Double); err != nil {
		return err
	}
	player := f.GetPlayer(pos)
	if f.doubled(pos) {
		return game.ErrCommandRepeated
	}
	f.state.Doubled = append(f.state.Doubled, int8(player.Pos()))
	if multiple > 1 {
		if err := player.Double(&f.state.Players[player.Pos()], multiple); err != nil {
			return err
		}
	}
	if len(f.state.Doubled) == f.NumPlayer() {
		f.state.Stage = Stage_Play
	}
	return nil
}

func (f *Landlord3) OnRevealCards(pos int) error {
	// check stage
	if err := f.expectStages(Stage_Play); err != nil {
		return err
	}
	if len(f.state.LastHands.LastPokers) > 0 {
		return game.ErrState
	}
	player := f.GetPlayer(pos)
	return player.RevealCards(&f.state.Players[player.Pos()])
}

func (f *Landlord3) OnPlayCards(pos int, pokers []uint8) error {
	if len(pokers) == 0 {
		return ErrEmptyPokers
	}
	return f.playCards(pos, pokers)
}

func (f *Landlord3) OnPass(pos int) error {
	return f.playCards(pos, nil)
}

func HandsOfPokers(pokers []uint8) Hands {
	//TODO
	return Hands{}
}

func (f *Landlord3) playCards(pos int, pokers []uint8) error {
	player := f.GetPlayer(pos)
	// check stage
	if err := f.expectStages(Stage_Play); err != nil {
		return err
	}
	// check turn
	if err := f.expectTurn(player.Pos()); err != nil {
		return err
	}
	if len(pokers) >= 0 {
		hands := HandsOfPokers(pokers)
		if !hands.Valid {
			return ErrInvalidCategory
		}
		lastHands := f.state.LastHands
		if int8(player.Pos()) != f.state.LastPos && lastHands.Valid {
			// check pokers
			isBomb := hands.Category == Category_Bomb || hands.Category == Category_Rocket
			lastIsBomb := lastHands.Category == Category_Bomb || lastHands.Category == Category_Rocket
			ok := false
			v1, v2 := hands.Value, lastHands.Value
			if isBomb {
				ok = !lastIsBomb || v1 > v2
			} else if !lastIsBomb {
				ok = hands.Category == lastHands.Category && v1 > v2
			}
			if !ok {
				return ErrInvalidCategory
			}
		}
		if err := player.PlayCards(&f.state.Players[player.Pos()], pokers); err != nil {
			return err
		}
		// TODO: double for bomb or rocket
		f.state.LastPos = int8(player.Pos())
		f.state.LastHands = hands
		f.storeLastPokers(pokers)
		if f.gameover() {
			return f.beforeGameover(false)
		}
	}
	f.nextTurn()
	return nil
}
