package orm

type CoreEngine interface {
	Name() string
	SetErrorHandler(ErrorHandler)
	CreateIndex(tableName string, index Index)
	Insert(table ReadonlyTable) error
	Update(table ReadonlyTable, fields ...string) error
	Find(meta TableMeta, keys KeyList, fs FieldSetterFactory, fields ...string) error
	Get(table WriteonlyTable, fields ...string) (bool, error)
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
			return "table `" + meta.Name() + "` GetField `" + field + "`", ErrFieldNotFound
		}
		args = append(args, value)
	}
	_, err := eng.redisc.HsetMulti(args...)
	return "redis.Hmset", err
}

func (eng *coreEngine) remove(tableName string, tableKey interface{}, fields []string) (string, error) {
	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, eng.tableName(tableName))
	for _, field := range fields {
		args = append(args, eng.fieldName(tableKey, field))
	}
	_, err := eng.redisc.Hdel(args...)
	return "redis.Hdel", err
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
		return "redis.hmget", false, err
	}
	if len(values) != fieldSize {
		return "redis.hmget", false, ErrUnexpectedLength
	}
	found := false
	for i := 0; i < fieldSize; i++ {
		if values[i] != nil {
			if err := table.SetField(fields[i], *values[i]); err != nil {
				return "table `" + meta.Name() + "` SetField `" + fields[i] + "`", false, err
			}
			found = true
		}
	}
	return "", found, nil
}

func (eng *coreEngine) find(meta TableMeta, keys KeyList, fs FieldSetterFactory, fields ...string) (string, error) {
	keySize := keys.Len()
	if keySize == 0 {
		return "", nil
	}
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	fieldSize := len(fields)
	args := make([]interface{}, 0, fieldSize*keySize+1)
	args = append(args, eng.tableName(meta.Name()))
	for i := 0; i < keySize; i++ {
		key := ToString(keys.Key(i))
		for _, field := range fields {
			args = append(args, eng.fieldName(key, field))
		}
	}
	_, values, err := eng.redisc.Hmgetstrings(args...)
	if err != nil {
		return "redis.hmget", err
	}
	length := len(values)
	if length != fieldSize*keySize {
		return "redis.hmget", ErrUnexpectedLength
	}
	for i := 0; i+fieldSize <= length; i += fieldSize {
		index := i / fieldSize
		setter := fs.New(keys.Key(index))
		for j, field := range fields {
			if values[i+j] != nil {
				if err := setter.SetField(field, *values[i+j]); err != nil {
					return "table `" + meta.Name() + "` SetField `" + field + "`", err
				}
			}
		}
	}
	return "", nil
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
func (eng *coreEngine) Find(meta TableMeta, keys KeyList, fs FieldSetterFactory, fields ...string) error {
	action, err := eng.find(meta, keys, fs, fields...)
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
