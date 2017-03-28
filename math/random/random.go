package random

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

var (
	digits         = []byte("0123456789")
	lowercaseChars = []byte("abcdefghijklmnopqrstuvwxyz")
	uppercaseChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialChars   = []byte("~!@#$%^&*")
)

const (
	O_DIGIT = 1 << iota
	O_LOWER_CHAR
	O_UPPER_CHAR
	O_SPECIAL_CHAR
)

type cryptoSource struct{}

func (source cryptoSource) Seed(int64) {}
func (source cryptoSource) Int63() int64 {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		return DefaultSource.Int63()
	}
	b[0] &= 0x7F
	return int64(binary.BigEndian.Uint64(b[:]))
}

var (
	DefaultSource rand.Source = rand.NewSource(time.Now().UnixNano())
	CryptoSource  rand.Source = cryptoSource{}
)

func Int63(source rand.Source) int64 {
	if source == nil {
		source = DefaultSource
	}
	return source.Int63()
}

func Intn(n int, source rand.Source) int {
	return int(Int63(source) % int64(n))
}

func Bool(source rand.Source) bool {
	return Intn(2, source) == 1
}

func String(length int, source rand.Source, modes ...int) string {
	if length <= 0 {
		return ""
	}
	var mode int
	for _, m := range modes {
		mode |= m
	}
	if mode&(O_DIGIT|O_LOWER_CHAR|O_UPPER_CHAR|O_SPECIAL_CHAR) == 0 {
		mode = O_LOWER_CHAR | O_UPPER_CHAR
	}
	size := 0
	if mode&O_DIGIT != 0 {
		size += len(digits)
	}
	if mode&O_LOWER_CHAR != 0 {
		size += len(lowercaseChars)
	}
	if mode&O_UPPER_CHAR != 0 {
		size += len(uppercaseChars)
	}
	if mode&O_SPECIAL_CHAR != 0 {
		size += len(specialChars)
	}
	if source == nil {
		source = DefaultSource
	}
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		index := int(source.Int63() % int64(size))
		tmpSize := 0
		if mode&O_DIGIT != 0 {
			tmpSize = len(digits)
			if index < tmpSize {
				buf.WriteByte(digits[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_LOWER_CHAR != 0 {
			tmpSize = len(lowercaseChars)
			if index < tmpSize {
				buf.WriteByte(lowercaseChars[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_UPPER_CHAR != 0 {
			tmpSize = len(uppercaseChars)
			if index < tmpSize {
				buf.WriteByte(uppercaseChars[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_SPECIAL_CHAR != 0 {
			tmpSize = len(specialChars)
			if index < tmpSize {
				buf.WriteByte(specialChars[index])
				continue
			}
		}
		buf.WriteByte('-')
	}
	return buf.String()
}

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

type swapableSlice struct {
	swapper func(int, int)
	length  int
}

func (s swapableSlice) Len() int      { return s.length }
func (s swapableSlice) Swap(i, j int) { s.swapper(i, j) }

func ShuffleSlice(slice interface{}, source rand.Source) {
	Shuffle(swapableSlice{
		swapper: reflect.Swapper(slice),
		length:  reflect.ValueOf(slice).Len(),
	}, source)
}
