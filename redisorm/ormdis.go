package orm

type RedisClient interface {
	HsetMulti(...interface{}) (string, error)
	Hmgetstrings(...interface{}) (int, []*string, error)
	Hdel(...interface{}) (string, error)
}

type ErrorHandler func(action string, err error) error

const ErrorHnadlerDepth = 4

type Engine struct {
	*coreEngine
}

func NewEngine(name string, redisc RedisClient) *Engine {
	core := &coreEngine{
		name:    name,
		redisc:  redisc,
		indexes: map[string]map[string]Index{},
	}
	return &Engine{
		coreEngine: core,
	}
}

// Core returns CoreEngine
func (eng *Engine) Core() CoreEngine {
	return eng.coreEngine
}

// Insert inserts hash_table `table`
func (eng *Engine) Insert(table interface{}) error {
	return nil
}

// MultiInsert is multi-version of `Insert`
func (eng *Engine) MultiInsert(tables interface{}) error {
	return nil
}

// Update updates hash_table `table` specific `fields`,
// only updates non-zero fields if len(fields) is zero
func (eng *Engine) Update(table interface{}, fields ...interface{}) error {
	return nil
}

// MultiUpdate is multi-version of `Update`
func (eng *Engine) MultiUpdate(tables interface{}, fields ...interface{}) error {
	return nil
}

// UpdateAllFields updates all fields of hash_table `table`
func (eng *Engine) UpdateAllFields(table interface{}) error {
	return nil
}

// MultiUpdateAllFields is multi-version of `UpdateAllFields`
func (eng *Engine) MultiUpdateAllFields(tables interface{}) error {
	return nil
}

// Find gets records by keys, load the specific `fields` or all fields if len(fields) is zero
func (eng *Engine) Find(keys KeyList, results interface{}, fields ...string) error {
	return nil
}
