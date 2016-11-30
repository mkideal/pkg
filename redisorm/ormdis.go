package orm

type RedisClient interface {
	HsetMulti(...interface{}) (string, error)
	Hmgetstrings(...interface{}) (int, []*string, error)
	Hdel(...interface{}) (string, error)
}

type ErrorHandler func(action string, err error) error

const ErrorHandlerDepth = 4

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
	action := "reflect"
	rotable, err := reflectReadonlyTable(table)
	if err == nil {
		action, err = eng.update(rotable)
		if err == nil {
			return nil
		}
	}
	return eng.catch("Insert: "+action, err)
}

// MultiInsert is multi-version of `Insert`
func (eng *Engine) MultiInsert(tables interface{}) error {
	return nil
}

// Update updates hash_table `table` specific `fields`,
func (eng *Engine) Update(table interface{}, fields ...string) error {
	action := "reflect"
	rotable, err := reflectReadonlyTable(table)
	if err == nil {
		action, err = eng.update(rotable, fields...)
		if err == nil {
			return nil
		}
	}
	return eng.catch("Update: "+action, err)
}

// MultiUpdate is multi-version of `Update`
func (eng *Engine) MultiUpdate(tables interface{}, fields ...string) error {
	return nil
}

// Find gets records by keys, load the specific `fields` or all fields if len(fields) is zero
func (eng *Engine) Find(keys KeyList, results interface{}, fields ...string) error {
	return nil
}
