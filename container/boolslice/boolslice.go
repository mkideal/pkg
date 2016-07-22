package boolslice

type BoolSlice struct {
	data   []byte
	length int
}

func New() *BoolSlice {
	return &BoolSlice{data: []byte{}}
}

func NewWithCap(size int) *BoolSlice {
	return &BoolSlice{data: make([]byte, 0, (size>>3)+1), length: size}
}

func (s *BoolSlice) Len() int                     { return s.length }
func (s *BoolSlice) ij(index int) (i int, j byte) { return index >> 3, byte(index & 0x7) }

func (s *BoolSlice) Get(index int) bool {
	i, j := s.ij(index)
	return s.data[i]&(1<<j) != 0
}

func (s *BoolSlice) Set(index int, value bool) {
	i, j := s.ij(index)
	if value {
		s.data[i] |= 1 << j
	} else {
		s.data[i] &= ^(1 << j)
	}
}

func (s *BoolSlice) Push(value bool) {
	i, j := s.ij(s.length)
	if i >= len(s.data) {
		s.data = append(s.data, 0)
	}
	if value {
		s.data[i] |= 1 << j
	} else {
		s.data[i] &= ^(1 << j)
	}
	s.length++
}

func (s *BoolSlice) Pop() (value bool) {
	if s.length == 0 {
		panic("length == 0")
	}
	s.length--
	value = s.Get(s.length)
	if n := (s.length >> 3) + 1; n < len(s.data) {
		s.data = s.data[:n]
	}
	return
}

func (s *BoolSlice) Truncate(from, to int) {
	if to < s.length {
		s.length = to
		if n := (s.length >> 3) + 1; n < len(s.data) {
			s.data = s.data[:n]
		}
	}
	if from > 0 {
		l := to - from
		for i := 0; i < l; i++ {
			s.Set(i, s.Get(i+from))
		}
		s.length = l
		if n := (s.length >> 3) + 1; n < len(s.data) {
			s.data = s.data[:n]
		}
	}
}
