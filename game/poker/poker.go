package poker

import (
	"sort"

	"github.com/mkideal/pkg/math/random"
)

// kind of poker
const (
	Spade   = 1 // ♠️
	Heart   = 2 // ♥️
	Club    = 3 // ♣️
	Diamond = 4 // ♦️
	Joker   = 5
)

const (
	OrderOfColoredJoker = 14
	OrderOfPlainJoker   = 15
)

// 0 0 0 0 0 0 0 0
// |-------| |---|
//    (1)     (2)
//
// (1) represents order of poker (1,2,...,13,14,15)
// (2) represents kind of poker (1,2,3,4,5,6)
type Poker uint8

func New(kind int, order int) Poker {
	return Poker(kind | (order << 3))
}

func (c Poker) IsValid() bool {
	kind, order := c.Kind(), c.Order()
	return c.IsJoker() || (kind >= Spade && kind <= Diamond && order >= 1 && order <= 13)
}

func (c Poker) Order() int {
	return (int(c) & 0xF8) >> 3
}

func (c Poker) Kind() int {
	return int(c) & 0x7
}

func (c Poker) Value() int {
	order := c.Order()
	if order == 1 || order == 2 {
		// value of A is 14
		// value of 2 is 15
		return order + 13
	}
	return order
}

func (c Poker) IsJoker() bool {
	kind, order := c.Kind(), c.Order()
	return kind == Joker && (order == OrderOfColoredJoker || order == OrderOfPlainJoker)
}

// sort pokers by order
type ByOrder []uint8

func (by ByOrder) Len() int { return len(by) }
func (by ByOrder) Less(i, j int) bool {
	o1, o2 := Poker(by[i]).Order(), Poker(by[j]).Order()
	if o1 == o2 {
		return by[i] < by[j]
	}
	return o1 < o2
}
func (by ByOrder) Swap(i, j int) { by[i], by[j] = by[j], by[i] }

// sort pokers by value
type ByValue []uint8

func (by ByValue) Len() int { return len(by) }
func (by ByValue) Less(i, j int) bool {
	v1, v2 := Poker(by[i]).Value(), Poker(by[j]).Value()
	if v1 == v2 {
		return by[i] < by[j]
	}
	return v1 < v2
}
func (by ByValue) Swap(i, j int) { by[i], by[j] = by[j], by[i] }

// GetPokers returns 54 pokers
func GetPokers() []uint8 {
	return []uint8{
		// ♠️
		uint8(New(Spade, 1)),
		uint8(New(Spade, 2)),
		uint8(New(Spade, 3)),
		uint8(New(Spade, 4)),
		uint8(New(Spade, 5)),
		uint8(New(Spade, 6)),
		uint8(New(Spade, 7)),
		uint8(New(Spade, 8)),
		uint8(New(Spade, 9)),
		uint8(New(Spade, 10)),
		uint8(New(Spade, 11)),
		uint8(New(Spade, 12)),
		uint8(New(Spade, 13)),

		// ♥️
		uint8(New(Heart, 1)),
		uint8(New(Heart, 2)),
		uint8(New(Heart, 3)),
		uint8(New(Heart, 4)),
		uint8(New(Heart, 5)),
		uint8(New(Heart, 6)),
		uint8(New(Heart, 7)),
		uint8(New(Heart, 8)),
		uint8(New(Heart, 9)),
		uint8(New(Heart, 10)),
		uint8(New(Heart, 11)),
		uint8(New(Heart, 12)),
		uint8(New(Heart, 13)),

		// ♣️
		uint8(New(Club, 1)),
		uint8(New(Club, 2)),
		uint8(New(Club, 3)),
		uint8(New(Club, 4)),
		uint8(New(Club, 5)),
		uint8(New(Club, 6)),
		uint8(New(Club, 7)),
		uint8(New(Club, 8)),
		uint8(New(Club, 9)),
		uint8(New(Club, 10)),
		uint8(New(Club, 11)),
		uint8(New(Club, 12)),
		uint8(New(Club, 13)),

		// ♦️
		uint8(New(Diamond, 1)),
		uint8(New(Diamond, 2)),
		uint8(New(Diamond, 3)),
		uint8(New(Diamond, 4)),
		uint8(New(Diamond, 5)),
		uint8(New(Diamond, 6)),
		uint8(New(Diamond, 7)),
		uint8(New(Diamond, 8)),
		uint8(New(Diamond, 9)),
		uint8(New(Diamond, 10)),
		uint8(New(Diamond, 11)),
		uint8(New(Diamond, 12)),
		uint8(New(Diamond, 13)),

		uint8(New(Joker, OrderOfColoredJoker)),
		uint8(New(Joker, OrderOfPlainJoker)),
	}
}

// Deal deals pokers, pokers = GetPokers if len(pokers) is 0
// nums represents number of poker for each player
// source could be nil
func Deal(nums []int, pokers []uint8, source random.Source) (res [][]uint8, leftover []uint8) {
	if len(pokers) == 0 {
		pokers = GetPokers()
	}
	random.Shuffle(ByOrder(pokers), source)

	res = make([][]uint8, len(nums))
	usedTotalNum := 0
	for i, num := range nums {
		if usedTotalNum+num < len(pokers) {
			res[i] = pokers[usedTotalNum : usedTotalNum+num]
		} else {
			res[i] = pokers[usedTotalNum:]
		}
		sort.Sort(ByValue(res[i]))
		usedTotalNum += len(res[i])
	}
	leftover = pokers[usedTotalNum:]
	return
}
