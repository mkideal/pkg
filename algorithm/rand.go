package algorithm

import (
	"math/rand"
	"sort"
)

type SwapableSlice interface {
	Len() int
	Swap(i, j int)
}

func Shuffle(orders SwapableSlice, source rand.Source) {
	for i := orders.Len() - 1; i >= 0; i-- {
		if source == nil {
			orders.Swap(i, rand.Intn(i+1))
		} else {
			orders.Swap(i, int(source.Int63())%(i+1))
		}
	}
}

func ShuffleInts(orders []int, source rand.Source)       { Shuffle(sort.IntSlice(orders), source) }
func ShuffleFloats(orders []float64, source rand.Source) { Shuffle(sort.Float64Slice(orders), source) }
func ShuffleStrings(orders []string, source rand.Source) { Shuffle(sort.StringSlice(orders), source) }
