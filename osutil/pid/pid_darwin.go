// +build darwin

package pid

import (
	"syscall"
)

func isPidExist(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}
