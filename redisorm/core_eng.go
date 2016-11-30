package orm

type CoreEngine interface {
	Name() string
	SetErrorHandler(ErrorHandler)
	CreateIndex(tableName string, index Index)
	Insert(table ReadonlyTable) error
	Update(table ReadonlyTable, fields ...string) error
	Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error
	Get(table WriteonlyTable, fields ...string) (bool, error)
	Remove(meta TableMeta, tableKey interface{}) error
	RemoveByFields(tableName string, tableKey interface{}, fields []string) error
	LoadView(View, KeyList, FieldSetterList) error
}

const (
	action_null  = ""
	action_hmget = "redis.hmget"
	action_hdel  = "redis.hdel"
)

func action_get_field(table, field string) string {
	return "table `" + table + "` GetField `" + field + "`"
}
func action_set_field(table, field string) string {
	return "table `" + table + "` SetField `" + field + "`"
}

// coreEngine implements CoreEngine interface
type coreEngine struct {
	name    string
	redisc  RedisClient
	indexes map[string]map[string]Index

	errorHandler ErrorHandler
}

func (eng *coreEngine) catch(action string, err error) error {
	if err != nil && eng.errorHandler != nil {
		err = eng.errorHandler(action, err)
	}
	return err
}

func (eng *coreEngine) tableName(name string) string {
	return eng.name + "@" + name
}

func (eng *coreEngine) fieldName(key interface{}, field string) string {
	return ToString(key) + ":" + field
}

// update updates some fields of record
func (eng *coreEngine) update(table ReadonlyTable, fields ...string) (string, error) {
	var (
		meta = table.Meta()
		key  = table.Key()
	)
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	args := make([]interface{}, 0, len(fields)*2+1)
	args = append(args, eng.tableName(meta.Name()))
	for _, field := range fields {
		args = append(args, eng.fieldName(key, field))
		value, ok := table.GetField(field)
		if !ok {
			return action_get_field(meta.Name(), field), ErrFieldNotFound
		}
		args = append(args, value)
	}
	_, err := eng.redisc.HsetMulti(args...)
	return action_hmget, err
}

func (eng *coreEngine) remove(tableName string, tableKey interface{}, fields []string) (string, error) {
	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, eng.tableName(tableName))
	for _, field := range fields {
		args = append(args, eng.fieldName(tableKey, field))
	}
	_, err := eng.redisc.Hdel(args...)
	return action_hdel, err
}

func (eng *coreEngine) get(table WriteonlyTable, fields ...string) (string, bool, error) {
	meta := table.Meta()
	tableKey := ToString(table.Key())
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	fieldSize := len(fields)
	args := make([]interface{}, 0, fieldSize+1)
	args = append(args, eng.tableName(meta.Name()))
	for _, field := range fields {
		args = append(args, eng.fieldName(tableKey, field))
	}
	_, values, err := eng.redisc.Hmgetstrings(args...)
	if err != nil {
		return action_hmget, false, err
	}
	if len(values) != fieldSize {
		return action_hmget, false, ErrUnexpectedLength
	}
	found := false
	for i := 0; i < fieldSize; i++ {
		if values[i] != nil {
			if err := table.SetField(fields[i], *values[i]); err != nil {
				return action_set_field(meta.Name(), fields[i]), false, err
			}
			found = true
		}
	}
	return action_null, found, nil
}

func (eng *coreEngine) find(meta TableMeta, keys KeyList, setters FieldSetterList, fields []string) (string, error) {
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	_, action, err := eng.findByFields(meta.Name(), keys, setters, FieldSlice(fields), nil)
	return action, err
}

func (eng *coreEngine) findByFields(table string, keys KeyList, setters FieldSetterList, fields FieldList, refs map[string]View) (map[string]InterfaceKeys, string, error) {
	keySize := keys.Len()
	if keySize == 0 {
		return nil, action_null, nil
	}
	fieldSize := fields.Len()
	args := make([]interface{}, 0, fieldSize*keySize+1)
	args = append(args, eng.tableName(table))
	for i := 0; i < keySize; i++ {
		key := ToString(keys.Key(i))
		for i := 0; i < fieldSize; i++ {
			args = append(args, eng.fieldName(key, fields.Field(i)))
		}
	}
	_, values, err := eng.redisc.Hmgetstrings(args...)
	if err != nil {
		return nil, action_null, err
	}
	length := len(values)
	if length != fieldSize*keySize {
		return nil, action_null, ErrUnexpectedLength
	}
	var keysGroup map[string]InterfaceKeys
	if len(refs) > 0 {
		keysGroup = make(map[string]InterfaceKeys)
		for field := range refs {
			keysGroup[field] = InterfaceKeys(make([]interface{}, keySize))
		}
	}
	for i := 0; i+fieldSize <= length; i += fieldSize {
		index := i / fieldSize
		setter := setters.New(table, index, keys.Key(index))
		for j := 0; j < fieldSize; j++ {
			field := fields.Field(j)
			value := values[i+j]
			if value != nil {
				if err := setter.SetField(field, *value); err != nil {
					return nil, action_set_field(table, field), err
				}
			}
			if ks, ok := keysGroup[field]; ok {
				if value == nil {
					ks[index] = nil
				} else {
					ks[index] = *value
				}
				keysGroup[field] = ks
			}
		}
	}
	return nil, action_null, nil
}

//------------------
// core engine APIs
//------------------

// Name returns database name
func (eng *coreEngine) Name() string { return eng.name }

// SetErrorHandler sets handler for handling error
func (eng *coreEngine) SetErrorHandler(eh ErrorHandler) {
	eng.errorHandler = eh
}

// CreateIndex creates an index for specific table
func (eng *coreEngine) CreateIndex(tableName string, index Index) {
	idx, ok := eng.indexes[tableName]
	if !ok {
		idx = make(map[string]Index)
		eng.indexes[tableName] = idx
	}
	idx[index.Name()] = index
}

// Insert inserts a new record or updates all fields of record
func (eng *coreEngine) Insert(table ReadonlyTable) error {
	action, err := eng.update(table)
	if err == nil {
		return nil
	}
	return eng.catch("Insert: "+action, err)
}

// Update updates specific fields of record
func (eng *coreEngine) Update(table ReadonlyTable, fields ...string) error {
	action, err := eng.update(table, fields...)
	if err == nil {
		return nil
	}
	return eng.catch("Update: "+action, err)
}

// Find gets many records
func (eng *coreEngine) Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error {
	action, err := eng.find(meta, keys, setters, fields)
	if err == nil {
		return nil
	}
	return eng.catch("Find: "+action, err)
}

// Get gets one record by specific fields. It will gets all fields if fields is empty
func (eng *coreEngine) Get(table WriteonlyTable, fields ...string) (bool, error) {
	action, ok, err := eng.get(table, fields...)
	if err == nil {
		return ok, nil
	}
	return ok, eng.catch("Get: "+action, err)
}

// Remove removes one record
func (eng *coreEngine) Remove(meta TableMeta, tableKey interface{}) error {
	action, err := eng.remove(meta.Name(), tableKey, meta.Fields())
	if err == nil {
		return nil
	}
	return eng.catch("Remove: "+action, err)
}

// RemoveByFields removes some fields of table
func (eng *coreEngine) RemoveByFields(tableName string, tableKey interface{}, fields []string) error {
	action, err := eng.remove(tableName, tableKey, fields)
	if err == nil {
		return nil
	}
	return eng.catch("RemoveByFields: "+action, err)
}

// LoadView loads view by keys and store loaded data to setters
func (eng *coreEngine) LoadView(view View, keys KeyList, setters FieldSetterList) error {
	action, err := eng.recursivelyLoadView(view, keys, setters)
	if err == nil {
		return nil
	}
	return eng.catch("LoadView: "+action, err)
}
