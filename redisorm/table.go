package orm

// KeyList holds n keys
type KeyList interface {
	Len() int
	Key(int) interface{}
}

type IntKeys []int
type Int64Keys []int64
type Uint64Keys []uint64
type StringKeys []string
type InterfaceKeys []interface{}

func (keys IntKeys) Len() int       { return len(keys) }
func (keys Int64Keys) Len() int     { return len(keys) }
func (keys Uint64Keys) Len() int    { return len(keys) }
func (keys StringKeys) Len() int    { return len(keys) }
func (keys InterfaceKeys) Len() int { return len(keys) }

func (keys IntKeys) Key(i int) interface{}       { return keys[i] }
func (keys Int64Keys) Key(i int) interface{}     { return keys[i] }
func (keys Uint64Keys) Key(i int) interface{}    { return keys[i] }
func (keys StringKeys) Key(i int) interface{}    { return keys[i] }
func (keys InterfaceKeys) Key(i int) interface{} { return keys[i] }

// FieldList holds n fields
type FieldList interface {
	Len() int
	Field(int) string
}

// Field implements FieldList which atmost contains one value
type Field string

func (f Field) Len() int {
	if f == "" {
		return 0
	}
	return 1
}

func (f Field) Field(i int) string { return string(f) }

// Fields implements FieldList
type FieldSlice []string

func (fs FieldSlice) Len() int           { return len(fs) }
func (fs FieldSlice) Field(i int) string { return fs[i] }

//-----------------
// Basic interface
//-----------------

type FieldGetter interface {
	GetField(field string) (interface{}, bool)
}

type FieldSetter interface {
	SetField(field, value string) error
}

// TableMeta holds table meta information
type TableMeta interface {
	Name() string
	Fields() []string
}

//-------------------
// Compose interface
//-------------------

type TableInfo interface {
	Meta() TableMeta
	Key() interface{}
}

type ReadonlyTable interface {
	TableInfo
	FieldGetter
}

type WriteonlyTable interface {
	TableInfo
	FieldSetter
}

// Table represents a hash_table in redis
type Table interface {
	TableInfo
	FieldGetter
	FieldSetter
}

type FieldSetterList interface {
	New(table string, index int, key interface{}) FieldSetter
}
