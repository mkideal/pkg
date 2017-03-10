package poker

import (
	"sort"

	"github.com/mkideal/pkg/algorithm"
)

// 花色种类
const (
	Spade   = 1 // 黑桃 ♠️
	Heart   = 2 // 红桃 ♥️
	Club    = 3 // 梅花 ♣️
	Diamond = 4 // 方块 ♦️
	Joker1  = 5 // 小王
	Joker2  = 6 // 大王
)

// 最低3位表示花色种类
// 然后4位表示面值(0,13)
type Poker uint8

func New(kind int, order int) Poker {
	return Poker(kind | (order << 3))
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
		return order + 14
	}
	return order
}

type ByOrder []Poker

func (by ByOrder) Len() int           { return len(by) }
func (by ByOrder) Less(i, j int) bool { return by[i].Order() < by[j].Order() }
func (by ByOrder) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

type ByValue []Poker

func (by ByValue) Len() int           { return len(by) }
func (by ByValue) Less(i, j int) bool { return by[i].Value() < by[j].Value() }
func (by ByValue) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

var pokers = [...]Poker{
	New(Spade, 1),
	New(Spade, 2),
	New(Spade, 3),
	New(Spade, 4),
	New(Spade, 5),
	New(Spade, 6),
	New(Spade, 7),
	New(Spade, 8),
	New(Spade, 9),
	New(Spade, 10),
	New(Spade, 11),
	New(Spade, 12),
	New(Spade, 13),

	New(Heart, 1),
	New(Heart, 2),
	New(Heart, 3),
	New(Heart, 4),
	New(Heart, 5),
	New(Heart, 6),
	New(Heart, 7),
	New(Heart, 8),
	New(Heart, 9),
	New(Heart, 10),
	New(Heart, 11),
	New(Heart, 12),
	New(Heart, 13),

	New(Club, 1),
	New(Club, 2),
	New(Club, 3),
	New(Club, 4),
	New(Club, 5),
	New(Club, 6),
	New(Club, 7),
	New(Club, 8),
	New(Club, 9),
	New(Club, 10),
	New(Club, 11),
	New(Club, 12),
	New(Club, 13),

	New(Diamond, 1),
	New(Diamond, 2),
	New(Diamond, 3),
	New(Diamond, 4),
	New(Diamond, 5),
	New(Diamond, 6),
	New(Diamond, 7),
	New(Diamond, 8),
	New(Diamond, 9),
	New(Diamond, 10),
	New(Diamond, 11),
	New(Diamond, 12),
	New(Diamond, 13),

	New(Joker1, OrderOfJoker1),
	New(Joker2, OrderOfJoker2),
}

const (
	OrderOfJoker1 = 101
	OrderOfJoker2 = 102
)

func GetPokers(indexes []int) []Poker {
	res := make([]Poker, len(indexes))
	for i := range res {
		res[i] = pokers[indexes[i]]
	}
	return res
}

func newOrders(n int) []int {
	res := make([]int, n)
	for i := range res {
		res[i] = i
	}
	return res
}

// 发牌
func Deal(nums []int) (res [][]Poker, remains []Poker) {
	orders := newOrders(len(pokers))
	algorithm.ShuffleInts(orders, nil)

	res = make([][]Poker, len(nums))
	usedTotalNum := 0
	for i, num := range nums {
		if usedTotalNum+num <= len(orders) {
			res[i] = GetPokers(orders[usedTotalNum : usedTotalNum+num])
		} else {
			res[i] = GetPokers(orders[usedTotalNum:])
		}
		sort.Sort(ByValue(res[i]))
		usedTotalNum += len(res[i])
	}
	remains = GetPokers(orders[usedTotalNum:])
	return
}
