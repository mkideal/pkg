package static

type BoolOrError interface{}

// Assert check value
func Assert(value BoolOrError, errormsg string) {}

// TypeAssert check type of value
func TypeAssert(value interface{}, expected interface{}) {}
