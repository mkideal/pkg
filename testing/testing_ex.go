package testing

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

func (c *common) LogCalldepth(calldepth int, args ...interface{}) {
	c.logWithCalldepth(calldepth, fmt.Sprintln(args...))
}

func (c *common) logWithCalldepth(calldepth int, s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.output = append(c.output, decorateWithCalldepth(calldepth, s)...)
}

func decorateWithCalldepth(calldepth int, s string) string {
	_, file, line, ok := runtime.Caller(calldepth) // decorate + log + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}
