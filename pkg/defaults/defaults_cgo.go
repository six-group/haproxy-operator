//go:build cgo

package defaults

// #include "../../include/haproxy/defaults.h"
import "C"

var MaxLineArgs = C.MAX_LINE_ARGS
