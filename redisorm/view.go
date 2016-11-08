package orm

type View interface {
	Table() string
	Fields() FieldList
	Refs() map[string]View
}

func (eng *coreEngine) LoadView(view View, keys KeyList, setters FieldSetterList) error {
	return eng.recursivelyLoadView(view, keys, setters)
}

func (eng *coreEngine) recursivelyLoadView(view View, keys KeyList, setters FieldSetterList) error {
	keysGroup, err := eng.loadView(view, keys, setters)
	if err != nil {
		return err
	}
	refs := view.Refs()
	if refs == nil {
		return nil
	}
	if len(keysGroup) != len(refs) {
		return ErrUnexpectedLength
	}
	for field, ref := range refs {
		if tmpKeys, ok := keysGroup[field]; ok {
			if err := eng.recursivelyLoadView(ref, tmpKeys, setters); err != nil {
				return err
			}
		} else {
			return ErrViewRefFieldMissing
		}
	}
	return nil
}

func (eng *coreEngine) loadView(view View, keys KeyList, setters FieldSetterList) (map[string]KeyList, error) {
	eng.findByFields(view.Table(), keys, setters, view.Fields(), view.Refs())
	return nil, nil
}
