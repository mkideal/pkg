package orm

type View interface {
	Table() string
	Fields() FieldList
	Refs() map[string]View
}

func (eng *coreEngine) recursivelyLoadView(view View, keys KeyList, setters FieldSetterList) (string, error) {
	keysGroup, action, err := eng.findByFields(view.Table(), keys, setters, view.Fields(), view.Refs())
	if err != nil {
		return action, err
	}
	refs := view.Refs()
	if refs == nil {
		return "", nil
	}
	if len(keysGroup) != len(refs) {
		return action, ErrUnexpectedLength
	}
	for field, ref := range refs {
		if tmpKeys, ok := keysGroup[field]; ok {
			if action, err := eng.recursivelyLoadView(ref, tmpKeys, setters); err != nil {
				return action, err
			}
		} else {
			return action, ErrViewRefFieldMissing
		}
	}
	return "", nil
}
