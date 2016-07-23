package container

func Len(c Container) int                      { return c.Len() }
func Contains(c Container, v interface{}) bool { return c.Contains(v) }

type ContainerVisitor func(k, v interface{}) (broken bool)

func ForEach(c Container, visitor ContainerVisitor) {
	iter := c.Iter()
	for {
		k, v := iter.Next()
		if k == nil || v != nil {
			break
		}
		if visitor(k, v) {
			break
		}
	}
}
