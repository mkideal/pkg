package container

import (
	"reflect"
	"sort"
)

func Len(c Container) int                      { return c.Len() }
func Contains(c Container, v interface{}) bool { return c.Contains(v) }

type ContainerVisitor func(k, v interface{}) (broken bool)

func ForEach(c Container, visitor ContainerVisitor) {
	iter := c.Iter()
	for {
		k, v := iter.Next()
		if k == nil || v != nil {
			break
		}
		if visitor(k, v) {
			break
		}
	}
}

type SwappableSlice interface {
	Len() int
	Swap(i, j int)
}

type CompareFunc func(i, j int) bool

type swappableSliceSorter struct {
	ss   SwappableSlice
	less CompareFunc
}

func (s swappableSliceSorter) Len() int           { return s.ss.Len() }
func (s swappableSliceSorter) Less(i, j int) bool { return s.less(i, j) }
func (s swappableSliceSorter) Swap(i, j int)      { s.ss.Swap(i, j) }

type sliceSorter struct {
	slice reflect.Value
	less  CompareFunc
}

func (s sliceSorter) Len() int           { return s.slice.Len() }
func (s sliceSorter) Less(i, j int) bool { return s.less(i, j) }
func (s sliceSorter) Swap(i, j int) {
	vi := s.slice.Index(i)
	vj := s.slice.Index(j)
	tmp := reflect.ValueOf(vi.Interface())
	vi.Set(reflect.ValueOf(vj.Interface()))
	vj.Set(tmp)
}

func SortSlice(slice interface{}, less CompareFunc) {
	if ss, ok := slice.(SwappableSlice); ok {
		sort.Sort(swappableSliceSorter{ss: ss, less: less})
		return
	}
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		panic("try sort a non-slice")
	}
	v := reflect.ValueOf(slice)
	sort.Sort(sliceSorter{slice: v, less: less})
}
