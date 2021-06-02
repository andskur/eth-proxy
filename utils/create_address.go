package utils

import (
	"fmt"
)

// CreateAddr is create address string from host and port
func CreateAddr(host string, port int) string {
	return fmt.Sprintf("%s:%v", host, port)
}
