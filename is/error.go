package is

type Box struct {
	err error
}

func NewBox() *Box {
	return &Box{}
}

func (box Box) Err() error {
	return box.err
}

func Nil(box *Box, err error) bool {
	box.err = err
	return err == nil
}

func Err(box *Box, err error) bool {
	box.err = err
	return err != nil
}
