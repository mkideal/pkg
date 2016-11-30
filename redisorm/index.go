package orm

type Index interface {
	Name() string
	RefTable() string
	Update(member interface{}, score int64) error
}
